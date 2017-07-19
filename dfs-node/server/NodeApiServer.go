package server

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"pilicat-core/config"
	"pilicat-core/logs"
	"pilicat-core/utils"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

type NodeApiServer struct {
}

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

func (_self *NodeApiServer) FileHandler(ctx context.Context) {
	fileName, data, err := _self.ParseMultipartUpload(ctx.Request())
	if err == nil {
		fileExt := utils.FileExt(fileName)
		fmt.Println("FileName: ", fileName, " FileExt: ", fileExt)
		filePath := fmt.Sprint("./file/", _self.RandomFileName(), fileExt)
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		defer f.Close()
		if err != nil {
			logs.Error("OpenFile error:", err)
		}
		//
		io.Copy(f, bytes.NewReader(data))

	} else {
		logs.Error("Upload file error:", err)
	}
	fmt.Println(time.Now())
	return
}

func (_self *NodeApiServer) MonthDirName() string {
	return utils.FormatTimeByFm(time.Now(), "200601")
}

func (_self *NodeApiServer) RandomFileName() string {
	return utils.GetGuid()
}

func (_self *NodeApiServer) DefaultHandler(ctx context.Context) {

	fmt.Println("DefaultHandler ", ctx.Method())

	OutputText(ctx, fmt.Sprint("Pilicat Dfs-Node Api ", Version))
}

func (_self *NodeApiServer) Run() {
	endRunning := make(chan bool, 1)
	nodeApiAddr := config.AppConf.DefaultString("node.api.addr", ":8800")

	logs.Info("Start dfs-node api server", nodeApiAddr)

	app := iris.New()
	app.Post("/api/file", context.LimitRequestBodySize(32<<20), _self.FileHandler)
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
