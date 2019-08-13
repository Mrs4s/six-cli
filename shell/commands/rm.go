package commands

import (
	"fmt"
	pl "github.com/Mrs4s/power-liner"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/shell"
	"strings"
)

func init() {
	alias["Delete"] = []string{"rm", "del"}
	explains["Delete"] = "删除文件或目录"
}

func (CommandHandler) Delete(c *pl.Context) {
	targets := models.FilterStrings(c.Nokeys, func(s string) bool { return s != "" && s != " " })
	for _, path := range targets {
		if _, ok := c.Keys["y"]; !ok {
			if in, _ := shell.App.ReadLine("[?] 确认是否删除 " + models.GetFileName(path) + " (y/n):"); strings.ToLower(in) != "y" {
				continue
			}
		}
		if strings.HasPrefix(path, "/") {
			if err := shell.CurrentUser.DeleteFile(path); err != nil {
				fmt.Println("[!] 文件", path, "删除失败:", err)
			}
			continue
		}
		t := models.CombinePaths(shell.CurrentPath, path, "/")
		if err := shell.CurrentUser.DeleteFile(t); err != nil {
			fmt.Println("[!] 文件", t, "删除失败:", err)
		}
	}
}

func (CommandHandler) DeleteCompleter(c *pl.Context) []string {
	if strings.HasSuffix(c.RawLine, " ") {
		return models.SelectStrings(append(filterCurrentDirs(), filterCurrentFiles()...), func(s string) string {
			if strings.Contains(s, " ") {
				return "\"" + s + "\""
			}
			return s
		})
	}
	return []string{}
}
