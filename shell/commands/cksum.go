package commands

import (
	"fmt"
	pl "github.com/Mrs4s/power-liner"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/models/fs"
	"github.com/Mrs4s/six-cli/shell"
	"github.com/Mrs4s/six-cli/six_cloud"
	"strconv"
	"strings"
)

func init() {
	alias["CheckSum"] = []string{"cksum"}
	explains["CheckSum"] = "效验文件hash"
}

func (CommandHandler) CheckSum(c *pl.Context) {
	targets := models.FilterStrings(c.Nokeys, func(s string) bool { return s != "" && s != " " })
	if len(targets) == 0 {
		fmt.Println("[H] 使用方法: cksum <文件...>")
		return
	}
	table := [][]string{{"Hash", "文件大小(字节)", "文件名"}}
	for _, file := range targets {
		files, err := shell.CurrentUser.GetFilesByPath(shell.CurrentPath)
		if strings.HasPrefix(file, "/") {
			files, err = shell.CurrentUser.GetFilesByPath(fs.GetParentPath(file))
		}
		if err != nil {
			fmt.Println("[!] 获取文件", fs.GetFileName(file), "信息失败:", err)
			continue
		}
		var target *six_cloud.SixFile
		for _, sub := range files {
			if sub.Name == fs.GetFileName(file) {
				target = sub
				break
			}
		}
		if target == nil {
			fmt.Println("[!] 获取文件", fs.GetFileName(file), "信息失败: 文件不存在")
			continue
		}
		if target.IsDir {
			fmt.Println("[!] 获取文件", fs.GetFileName(file), "信息失败: 目标为文件夹")
			continue
		}
		table = append(table, []string{target.ETag, strconv.FormatInt(target.Size, 10), target.Name})
	}
	shell.App.PrintTables(table, 2)
}

func (CommandHandler) CheckSumCompleter(c *pl.Context) []string {
	if strings.HasSuffix(c.RawLine, " ") {
		return models.SelectStrings(filterCurrentFiles(), func(s string) string {
			if strings.Contains(s, " ") {
				return "\"" + s + "\""
			}
			return s
		})
	}
	return []string{}
}
