package config

import (
	"os"
	"path/filepath"
)

var (
	AppConf Configer
)

func init() {

	loadAppConf()

}

func loadAppConf() {
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	appConfigPath := filepath.Join(workPath, "conf", "app.conf")

	AppConf, err = NewConfig("ini", appConfigPath)
	if err != nil {
		panic(err)
	}
}
