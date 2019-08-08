package main

import (
	"fmt"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/six_cloud"
	"os"
)

var (
	currentUser *six_cloud.SixUser
	currentPath = "/"
)

func main() {
	if !models.PathExists("config.json") {
		tmpConf := &models.Config{
			DownloadThread: 4, //Default
		}
		tmpConf.SaveFile("config.json")
	}
	models.DefaultConf = models.LoadConfig("config.json")
	if models.DefaultConf.QingzhenToken != "" {
		currentUser, _ = six_cloud.LoginWithAccessToken(models.DefaultConf.QingzhenToken)
	}
	if models.DefaultConf.DownloadPath == "" || !models.PathExists(models.DefaultConf.DownloadPath) {
		fmt.Println("下载路径不存在, 使用默认下载路径 <工作目录>/Downloads")
		models.DefaultConf.DownloadPath = "Downloads"
		if !models.PathExists("Downloads") {
			_ = os.MkdirAll("Downloads", os.ModePerm)
		}
	}
	runAsShell()
	/*
		if len(os.Args) == 1 {
			fmt.Println("usage: six-cli <command> or six-cli shell")
			return
		}
		if os.Args[1] == "shell" {
			runAsShell()
			return
		}
	*/
}
