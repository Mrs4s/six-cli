package main

import (
	"fmt"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/models/fs"
	"github.com/Mrs4s/six-cli/shell"
	"github.com/Mrs4s/six-cli/shell/commands"
	"os"
)

func main() {
	if !fs.PathExists("config.json") {
		tmpConf := &models.Config{
			DownloadThread:    4,
			DownloadBlockSize: 104857600,
			Mount: &models.MountConfig{
				AllowPreload: true,
				ChunkSize:    1024 * 1024,
				ChunkMax:     4,
			},
		}
		tmpConf.SaveFile("config.json")
	}
	models.DefaultConf = models.LoadConfig("config.json")
	if models.DefaultConf.DownloadPath == "" || !fs.PathExists(models.DefaultConf.DownloadPath) {
		fmt.Println("下载路径不存在, 使用默认下载路径 <工作目录>/Downloads")
		models.DefaultConf.DownloadPath = "Downloads"
		if !fs.PathExists("Downloads") {
			_ = os.MkdirAll("Downloads", os.ModePerm)
		}
	}
	if len(os.Args) == 1 {
		shell.RunAsShell(&commands.CommandHandler{})
		os.Exit(0)
	}
	shell.RunAsCli(&commands.CommandHandler{})
}
