package commands

import (
	"fmt"
	pl "github.com/Mrs4s/power-liner"
	"github.com/Mrs4s/six-cli/models"
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
		shell.CurrentPath = models.GetParentPath(shell.CurrentPath)
	case strings.Contains(arg, "../"):
		if shell.CurrentPath == "/" {
			return
		}
		for i := 0; i < strings.Count(arg, "../"); i++ {
			shell.CurrentPath = models.GetParentPath(shell.CurrentPath)
		}
	default:
		newPath := models.CombinePaths(shell.CurrentPath, arg, "/")
		_, err := shell.CurrentUser.GetFilesByPath(newPath)
		if err != nil {
			fmt.Println("[!] 切换失败: " + err.Error())
			return
		}
		shell.CurrentPath = newPath
	}
}

func (CommandHandler) JoinPathCompleter(c *pl.Context) []string {
	if len(c.Nokeys) > 1 {
		return []string{}
	}
	if len(strings.Split(c.Nokeys[0], "/")) <= 1 {
		return models.SelectStrings(filterCurrentDirs(), func(s string) string {
			if strings.Contains(s, " ") {
				return "\"" + s + "\""
			}
			return s
		})
	}
	newPath := shell.CurrentPath + "/" + c.Nokeys[0]
	if shell.CurrentPath == "/" {
		newPath = "/" + c.Nokeys[0]
	}
	files, err := shell.CurrentUser.GetFilesByPath(models.GetParentPath(newPath))
	if err != nil {
		return []string{}
	}
	return models.SelectStrings(models.FilterStrings(filterDirs(files), func(s string) bool {
		if strings.HasSuffix(newPath, "/") {
			return true
		}
		return strings.HasPrefix(s, models.GetFileName(newPath))
	}), func(s string) string {
		com := models.CombinePaths(models.GetParentPath(newPath), s, "/")
		if strings.Contains(com, " ") {
			return "\"" + com[1:] + "\""
		}
		return com[1:]
	})
}
