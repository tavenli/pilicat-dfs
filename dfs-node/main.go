package main

import (
	"fmt"
	"net/http"
	"os"
	"pilicat-core/config"
	"pilicat-core/logs"
	"time"

	"github.com/chrislusf/seaweedfs/weed/util"
)

func main() {

	logs.SetLogger(logs.AdapterConsole, `{"level":7}`)
	logs.SetLogger(logs.AdapterFile, `{"filename":"app.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)

	//输出文件名和行号
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	//为了让日志输出不影响性能，开启异步日志
	logs.Async()

	appName := config.AppConf.String("app.name")

	logs.Debug("★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★")
	logs.Debug("               ", appName, " 启动")
	logs.Debug("")
	logs.Debug("★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★")

	RunNode()
}

func RunNode() {

	nodeApiAddr := config.AppConf.DefaultString("node.api.addr", ":8800")

	endRunning := make(chan bool, 1)

	logs.Info("Start dfs-node api server", nodeApiAddr)
	listener, err := util.NewListener(nodeApiAddr, time.Duration(30)*time.Second)
	if err != nil {
		logs.Error("dfs-node api server listener error:", err)
	}

	apiNodeMux := initApiNodeMux()

	go func() {

		if err := http.Serve(listener, apiNodeMux); err != nil {
			logs.Critical("Dfs-Node: ", err, fmt.Sprintf("%d", os.Getpid()))
			time.Sleep(100 * time.Microsecond)
			endRunning <- true
		}
	}()

	<-endRunning

	logs.Info("Close dfs-node api server", nodeApiAddr)
	return
}

func initApiNodeMux() *http.ServeMux {
	apiNodeMux := http.NewServeMux()
	apiNodeMux.HandleFunc("/api/post", PostHandler)

	return apiNodeMux
}

func PostHandler(w http.ResponseWriter, r *http.Request) {

}
