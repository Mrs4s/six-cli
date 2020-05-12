package commands

import (
	"fmt"
	pl "github.com/Mrs4s/power-liner"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/shell"
	"github.com/Mrs4s/six-cli/six_cloud"
	"strconv"
	"time"
)

func init() {
	alias["Login"] = []string{}
	explains["Login"] = "登录6Pan账号"
}

func (CommandHandler) Login(c *pl.Context) {
	//args := models.FilterStrings(c.RawArgs, func(s string) bool { return strings.TrimSpace(s) != "" })
	token, _, err := six_cloud.CreateDestination()
	if err != nil {
		fmt.Println("[!] 创建网页Token失败, 请重试.")
		return
	}
	state := strconv.FormatInt(time.Now().Unix(), 10)
	fmt.Println("[+] 请在浏览器中打开以下链接, 并完成登录操作.")
	fmt.Println()
	fmt.Println(fmt.Sprintf("https://account.6pan.cn/login?destination=%s&appid=3cnu7s71h92p&response=query&state=%s&lang=zh-CN", token, state))
	fmt.Println()
	var user *six_cloud.SixUser
	for {
		u, err := six_cloud.LoginWithWebToken(token, state)
		if err != nil && err != six_cloud.ErrWaitingLogin {
			fmt.Println("[!] 登录失败")
			return
		}
		if u != nil {
			user = u
			break
		}
		time.Sleep(time.Second)
	}
	shell.CurrentUser = user
	shell.CurrentPath = "/"
	var flag bool
	for _, su := range shell.SavedUsers {
		if su.Identity == user.Identity {
			flag = true
		}
	}
	if !flag {
		models.DefaultConf.Tokens = append(models.DefaultConf.Tokens, user.Client.QingzhenToken)
		shell.SavedUsers = append(shell.SavedUsers, user)
	}
	fmt.Println("[+] 登录完成, 欢迎: " + user.Username)
	fmt.Println()
	if len(shell.SavedUsers) > 1 {
		printUserList()
	}
	models.DefaultConf.SaveFile("config.json")
	refreshPrompt()
}
