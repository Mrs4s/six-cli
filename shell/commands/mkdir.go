package commands

import (
	"fmt"
	pl "github.com/Mrs4s/power-liner"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/shell"
	"strings"
)

func init() {
	alias["Mkdir"] = []string{}
	explains["Mkdir"] = "创建目录"
}

func (CommandHandler) Mkdir(c *pl.Context) {
	targets := models.FilterStrings(c.Nokeys, func(s string) bool { return s != "" && s != " " })
	for _, path := range targets {
		if strings.HasPrefix(path, "/") {
			fmt.Println("[+] 创建目录", models.ShortPath(path, 30))
			if err := shell.CurrentUser.CreateDirectory(path); err != nil {
				fmt.Println("[!] 创建目录", path, "失败:", err)
			}
			continue
		}
		t := models.CombinePaths(shell.CurrentPath, path, "/")
		fmt.Println("[+] 创建目录", models.ShortPath(t, 30))
		if err := shell.CurrentUser.CreateDirectory(t); err != nil {
			fmt.Println("[!] 创建目录", t, "失败:", err)
		}
	}
	fmt.Println("[+] 操作完成")
}
