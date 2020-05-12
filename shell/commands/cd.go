package commands

import (
	"fmt"
	pl "github.com/Mrs4s/power-liner"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/models/fs"
	"github.com/Mrs4s/six-cli/shell"
	"strings"
)

func init() {
	alias["JoinPath"] = []string{"cd", "join"}
	explains["JoinPath"] = "切换工作目录"
}

func (CommandHandler) JoinPath(c *pl.Context) {
	if len(c.Nokeys) == 0 {
		fmt.Println("[H] 使用方法: cd <目录>")
	}
	arg := c.Nokeys[0]
	defer refreshPrompt()
	switch {
	case len(arg) == 0:
	case arg[0:1] == "/":
		_, err := shell.CurrentUser.GetFilesByPath(arg)
		if err != nil {
			fmt.Println("[!] 切换失败: " + err.Error())
			return
		}
		shell.CurrentPath = arg
	case arg == "..":
		if shell.CurrentPath == "/" {
			return
		}
		shell.CurrentPath = fs.GetParentPath(shell.CurrentPath)
	case strings.Contains(arg, "../"):
		if shell.CurrentPath == "/" {
			return
		}
		for i := 0; i < strings.Count(arg, "../"); i++ {
			shell.CurrentPath = fs.GetParentPath(shell.CurrentPath)
		}
	default:
		newPath := models.CombinePaths(shell.CurrentPath, arg, "/")
		if strings.HasSuffix(newPath, "/") {
			newPath = newPath[:len(newPath)-1]
		}
		_, err := shell.CurrentUser.GetFilesByPath(newPath)
		if err != nil {
			fmt.Println("[!] 切换失败: " + err.Error())
			return
		}
		shell.CurrentPath = newPath
	}
}

func (CommandHandler) JoinPathCompleter(c *pl.Context) []string {
	return PathCompleter(c, false)
}
