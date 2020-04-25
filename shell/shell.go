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
	SavedUsers  []*six_cloud.SixUser
	CurrentPath = "/"
	App         *pl.Shell
)

func RunAsShell(handler pl.IHandler) {
	if models.DefaultConf.Tokens != nil && len(models.DefaultConf.Tokens) > 0 {
		CurrentUser, _ = six_cloud.LoginWithAccessToken(models.DefaultConf.Tokens[0])
		for _, token := range models.DefaultConf.Tokens {
			if user, err := six_cloud.LoginWithAccessToken(token); err == nil {
				SavedUsers = append(SavedUsers, user)
			}
		}
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
		models.DefaultConf.Tokens = []string{}
		for ind := range SavedUsers {
			models.DefaultConf.Tokens = append(models.DefaultConf.Tokens, SavedUsers[ind].Client.QingzhenToken)
		}
		models.DefaultConf.SaveFile("config.json")
		os.Exit(0)
	})
	if CurrentUser != nil {
		fmt.Println("自动登录成功, 欢迎您", CurrentUser.Username)
		_, _ = fmt.Fprintf(os.Stdout, "当前用量: %s / %s  (%.2f%%)\n",
			models.ConvertSizeString(CurrentUser.UsedSpace),
			models.ConvertSizeString(CurrentUser.TotalSpace),
			float64(CurrentUser.UsedSpace)/float64(CurrentUser.TotalSpace)*100)
		if len(SavedUsers) > 1 {
			for _, user := range SavedUsers {
				if user.Identity == CurrentUser.Identity {
					fmt.Println("->", user.Username)
					continue
				}
				fmt.Println(user.Username)
			}
		}
		App.SetPrompt(CurrentUser.Username + "@six-pan:/$ ")
	}
	App.AddHandler(handler)
	App.RunAsShell()
}

func RunAsCli(handler pl.IHandler) {
	if models.DefaultConf.Tokens != nil && len(models.DefaultConf.Tokens) > 0 {
		CurrentUser, _ = six_cloud.LoginWithAccessToken(models.DefaultConf.Tokens[0])
	}
	if CurrentUser == nil {
		fmt.Println("[!] 请先登录!")
		return
	}
	App = pl.NewApp()
	App.AddHandler(handler)
	App.RunAsCli(os.Args[1:])
}
