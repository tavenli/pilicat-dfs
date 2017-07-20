package server

import (
	"bytes"
	"dfs-common"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"pilicat-core/config"
	"pilicat-core/logs"
	"pilicat-core/models"
	"pilicat-core/utils"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

type NodeApiServer struct {
}

var (
	dfsPubUrl string
)

type GenFilePath func(fName string) (dirName, fileName string)

func (_self *NodeApiServer) ParseFormUpload(ctx context.Context, key string) (fileName string, file multipart.File, e error) {
	file, info, err := ctx.FormFile(key)
	if err == nil {
		return info.Filename, file, err
	}

	return "", nil, err

	//io.Copy(out, file)
}

func (_self *NodeApiServer) ParseMultipartUpload(r *http.Request) (fileName string, data []byte, e error) {

	rReader, err := r.MultipartReader()
	if err != nil {
		logs.Error("MultipartReader error:", err)
		return "", nil, err
	}

	part, err := rReader.NextPart()
	if err != nil {
		logs.Error("Reading Multi part error:", err)
		return "", nil, err
	}

	fileName = part.FileName()
	if utils.IsNotEmpty(fileName) {
		fileName = path.Base(fileName)
		data, err := ioutil.ReadAll(part)
		return fileName, data, err
	}

	//if the filename is empty string, do a search on the other multi-part items
	for utils.IsEmpty(fileName) {
		logs.Debug("filename is empty")
		part, err = rReader.NextPart()
		if err != nil {
			// no more or on error, just safely break
			logs.Debug("NextPart error:", err)
			break
		}

		fName := part.FileName()

		//found the first <file type> multi-part has filename
		if utils.IsNotEmpty(fName) {
			data, err := ioutil.ReadAll(part)
			if err != nil {
				logs.Error("Reading Content error:", err)
				return "", nil, err
			}

			fileName = path.Base(fName)
			return fileName, data, err
		}
	}

	return
}

func (_self *NodeApiServer) GenFilePathByDaily(fName string) (dirName, fileName string) {
	fileExt := utils.FileExt(fName)
	rndFileName := fmt.Sprint(_self.RandomFileName(), fileExt)

	return _self.DailyDirName(), rndFileName
}

func (_self *NodeApiServer) FileHandler(ctx context.Context) {
	logs.Debug(utils.GetIP(ctx.Request()), " ", ctx.Method(), " ", ctx.Path())

	result := &models.ResultData{}

	fileName, data, err := _self.ParseMultipartUpload(ctx.Request())
	if err == nil {
		logs.Debug("Upload FileName: ", fileName)
		dirName, genFileName := _self.GenFilePathByDaily(fileName)

		fileCtxPath := fmt.Sprint("/file/", dirName)
		pwd, _ := os.Getwd()
		fileDir := fmt.Sprint(pwd, fileCtxPath)
		err := utils.MakeDir(fileDir, os.ModePerm)
		if err != nil {
			logs.Error("MakeDir error:", err)
			OutputJson(ctx, models.ResultData{Code: 1, Msg: err.Error()})
			return
		}
		filePath := fmt.Sprint(fileDir, "/", genFileName)
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		defer f.Close()
		if err != nil {
			logs.Error("OpenFile error:", err)
			OutputJson(ctx, models.ResultData{Code: 1, Msg: err.Error()})
			return
		}
		//
		io.Copy(f, bytes.NewReader(data))
		urlCtxPath := fmt.Sprint(fileCtxPath, "/", genFileName)
		logs.Debug("Saved File. ", urlCtxPath)
		dfsFile := new(dfsCommon.DfsFileInfo)
		dfsFile.FileUrlPath = urlCtxPath
		dfsFile.OrgFileName = fileName
		dfsFile.PubUrl = fmt.Sprint(dfsPubUrl, urlCtxPath)

		result.Code = 0
		result.Msg = "success"
		result.Data = dfsFile

	} else {
		logs.Error("Upload file error:", err)
		result.Code = 1
		result.Msg = err.Error()
	}

	//OutputJson(ctx, models.ResultData{Code: 0, Msg: ""})
	OutputJson(ctx, result)
	return
}

func (_self *NodeApiServer) MonthDirName() string {
	return utils.FormatTimeByFm(time.Now(), "200601")
}

func (_self *NodeApiServer) DailyDirName() string {
	return utils.FormatTimeByFm(time.Now(), "2006/01/02")
}

func (_self *NodeApiServer) RandomFileName() string {
	return utils.GetGuid()
}

func (_self *NodeApiServer) DefaultHandler(ctx context.Context) {

	logs.Debug("DefaultHandler ", ctx.Method())

	OutputText(ctx, fmt.Sprint("Pilicat Dfs-Node Api ", Version))
}

func (_self *NodeApiServer) Run() {
	endRunning := make(chan bool, 1)

	nodeApiAddr := config.AppConf.DefaultString("node.api.addr", ":8800")
	dfsPubUrl = config.AppConf.DefaultString("dfs.public.url", "http://your.domain")

	logs.Info("Start dfs-node api server", nodeApiAddr)

	app := iris.New()
	//不限制上传文件的大小
	app.Post("/api/file", _self.FileHandler)
	//限制上传文件的大小，32<<20 表示 32M
	//app.Post("/api/file", context.LimitRequestBodySize(32<<20), _self.FileHandler)
	app.Any("/", _self.DefaultHandler)

	go func() {

		if err := app.Run(iris.Addr(nodeApiAddr)); err != nil {
			logs.Critical("Dfs-Node: ", err, fmt.Sprintf("%d", os.Getpid()))
			time.Sleep(100 * time.Microsecond)
			endRunning <- true
		}

	}()

	<-endRunning

	logs.Info("Close dfs-node api server", nodeApiAddr)

}
