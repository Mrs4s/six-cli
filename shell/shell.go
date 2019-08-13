package shell

import (
	"fmt"
	pl "github.com/Mrs4s/power-liner"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/six_cloud"
	"os"
)

var (
	CurrentUser *six_cloud.SixUser
	CurrentPath = "/"
	App         *pl.Shell
)

func RunAsShell(handler pl.IHandler) {
	if models.DefaultConf.QingzhenToken != "" {
		CurrentUser, _ = six_cloud.LoginWithAccessToken(models.DefaultConf.QingzhenToken)
	}
	App = pl.NewApp()
	fmt.Println("欢迎使用6 Pan命令行客户端!")
	fmt.Println("开发By Mrs4s")
	fmt.Println("使用前请先前往 https://github.com/Mrs4s/six-cli 阅读使用指南")
	fmt.Println()
	App.SetPrompt("guest@six-pan:/$ ")
	App.OnAbort(func(count int32) {
		if count < 2 {
			fmt.Println("再次键入 Ctrl+C 以确认退出")
			return
		}
		models.DefaultConf.SaveFile("config.json")
		os.Exit(0)
	})
	if CurrentUser != nil {
		fmt.Println("自动登录成功, 欢迎您", CurrentUser.Username)
		App.SetPrompt(CurrentUser.Username + "@six-pan:/$ ")
	}
	App.AddHandler(handler)
	App.RunAsShell()
}

func RunAsCli(handler pl.IHandler) {
	if models.DefaultConf.QingzhenToken != "" {
		CurrentUser, _ = six_cloud.LoginWithAccessToken(models.DefaultConf.QingzhenToken)
	}
	App = pl.NewApp()
	App.AddHandler(handler)
	App.RunAsCli(os.Args[1:])
}
