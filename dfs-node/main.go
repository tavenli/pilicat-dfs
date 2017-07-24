package main

import (
	"dfs-node/server"
	"pilicat-core/config"
	"pilicat-core/logs"
	"time"
)

var (
	nodePubServer *server.NodePubServer
	nodeApiServer *server.NodeApiServer
)

func main() {

	logs.SetLogger(logs.AdapterConsole, `{"level":7}`)
	logs.SetLogger(logs.AdapterFile, `{"filename":"app.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)

	//输出文件名和行号
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	//为了让日志输出不影响性能，开启异步日志
	logs.Async()

	nodeName := config.AppConf.String("node.name")

	logs.Debug("★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★")
	logs.Debug("               ", nodeName, " 启动")
	logs.Debug("")
	logs.Debug("★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★")

	nodePubServer.Run()
	nodeApiServer.Run()

	logs.Info("The service will stop.")
	time.Sleep(1 * time.Second)
}
