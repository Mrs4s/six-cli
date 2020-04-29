package commands

import (
	"fmt"
	pl "github.com/Mrs4s/power-liner"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/shell"
	"github.com/Mrs4s/six-cli/six_cloud"
	"strconv"
	"strings"
	"time"
)

func init() {
	alias["Find"] = []string{}
	explains["Find"] = "搜索文件"
}

func (CommandHandler) Find(c *pl.Context) {
	if len(c.Nokeys) != 1 {
		fmt.Println("使用方法: find 文件名 [-d] [-s]")
		return
	}
	var files []*six_cloud.SixFile
	if strings.Contains(c.Nokeys[0], "?") || strings.Contains(c.Nokeys[0], "*") {
		var blocks []string
		for _, str := range strings.Split(c.Nokeys[0], "?") {
			blocks = append(blocks, strings.Split(str, "*")...)
		}
		for _, block := range models.FilterStrings(blocks, func(s string) bool { return s != "" }) {
			tmp, err := shell.CurrentUser.SearchFilesByName("", block)
			if err == nil {
				for _, t := range tmp {
					var flag bool
					for _, file := range files {
						if file.Path == t.Path {
							flag = true
							break
						}
					}
					if !flag && models.ShellMatch(t.Name, c.Nokeys[0]) {
						files = append(files, t)
					}
				}
			}
		}
	} else {
		files, _ = shell.CurrentUser.SearchFilesByName("", c.Nokeys[0])
	}
	table := [][]string{{"序号", "创建时间", "文件大小", "文件路径"}}
	if _, ok := c.Keys["d"]; ok {
		files = filterFileArray(files, func(file *six_cloud.SixFile) bool { return file.IsDir })
	}
	if _, ok := c.Keys["s"]; ok {
		files = filterFileArray(files, func(file *six_cloud.SixFile) bool { return file.Shared })
	}
	for i, file := range files {
		table = append(table, []string{
			strconv.FormatInt(int64(i), 10),
			time.Unix(file.CreateTime/1000, 0).In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05"),
			models.ConvertSizeString(file.Size),
			file.Path,
		})
	}
	shell.App.PrintTables(table, 2)
}
