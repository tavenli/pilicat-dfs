package server

import (
	"fmt"
	"os"
	"pilicat-core/config"
	"pilicat-core/logs"
	"pilicat-core/utils"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

type NodePubServer struct {
}

func (_self *NodePubServer) DefaultHandler(ctx context.Context) {
	logs.Debug(utils.GetIP(ctx.Request()), " ", ctx.Method(), " ", ctx.Path())

	OutputText(ctx, fmt.Sprint("Pilicat Dfs-Node Pub ", Version))
}

func (_self *NodePubServer) Run() {
	nodePubAddr := config.AppConf.DefaultString("node.public.addr", "")

	if utils.IsNotEmpty(nodePubAddr) {
		logs.Info("Start dfs-node pub server", nodePubAddr)

		app := iris.New()
		//app.Favicon("./favicon.ico")
		app.StaticWeb("/file", "./file")
		app.Any("/", _self.DefaultHandler)

		go func() {

			if err := app.Run(iris.Addr(nodePubAddr)); err != nil {
				logs.Critical("Dfs-Node: ", err, fmt.Sprintf("%d", os.Getpid()))
				time.Sleep(100 * time.Microsecond)

			}
		}()

	}
}
