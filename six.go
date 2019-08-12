package main

import (
	"fmt"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/shell"
	"github.com/Mrs4s/six-cli/shell/commands"
	"os"
)

func main() {
	if !models.PathExists("config.json") {
		tmpConf := &models.Config{
			DownloadThread: 16,
			PeakTaskCount:  2,
		}
		tmpConf.SaveFile("config.json")
	}
	models.DefaultConf = models.LoadConfig("config.json")
	if models.DefaultConf.DownloadPath == "" || !models.PathExists(models.DefaultConf.DownloadPath) {
		fmt.Println("下载路径不存在, 使用默认下载路径 <工作目录>/Downloads")
		models.DefaultConf.DownloadPath = "Downloads"
		if !models.PathExists("Downloads") {
			_ = os.MkdirAll("Downloads", os.ModePerm)
		}
	}
	shell.RunAsShell(&commands.CommandHandler{})
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
