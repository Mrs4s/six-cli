package commands

import (
	"fmt"
	pl "github.com/Mrs4s/power-liner"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/shell"
	"strconv"
	"strings"
	"time"
)

func init() {
	explains["Offline"] = "离线下载操作"
}

//Offline command
func (CommandHandler) Offline(c *pl.Context) {
	if len(c.RawArgs) == 1 {
		fmt.Println("使用方法: offline <list/add/del/filter> [...]")
		return
	}
	switch strings.ToLower(c.Nokeys[0]) {
	case "list":
		tasks, err := shell.CurrentUser.GetOfflineTasks()
		if err != nil {
			fmt.Println("获取失败:", err)
			return
		}
		tables := [][]string{{"序号", "创建时间", "状态", "进度", "任务名"}}
		if _, ok := c.Keys["e"]; ok {
			tables[0] = []string{"序号", "状态", "错误信息", "任务名"}
			for i, task := range tasks {
				if task.ErrorCode != 0 {
					tables = append(tables, []string{
						strconv.FormatInt(int64(i), 10),
						task.StatusStr(),
						task.ErrorMessage,
						task.Name,
					})
				}
			}
			shell.App.PrintTables(tables, 2)
			return
		}
		for i, task := range tasks {
			tables = append(tables, []string{
				strconv.FormatInt(int64(i), 10),
				time.Unix(task.CreateTime/1000, 0).In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05"),
				task.StatusStr(),
				strconv.FormatInt(int64(task.Progress), 10) + "%",
				task.Name,
			})
		}
		shell.App.PrintTables(tables, 2)
	case "add":
		if len(c.Nokeys) == 1 {
			fmt.Println("使用方法: offline add <链接> [-p 密码] [-o 输出目录] [-y]")
			fmt.Println("注意: 不添加输出目录参数的话将默认添加到工作目录")
			return
		}
		url := c.Nokeys[1]
		pass := c.Keys["p"]
		fmt.Println("正在预解析....")
		identity, name, size, err := shell.CurrentUser.PreparseOffline(url, pass)
		if err != nil {
			fmt.Println("解析错误:", err)
			return
		}
		if _, ok := c.Keys["y"]; !ok {
			if i, _ := shell.App.ReadLine("解析完成, 任务名: " + name + " 任务大小: " + models.ConvertSizeString(size) + " 是否确定下载? (y/n)"); i != "y" {
				return
			}
		}
		fmt.Println("正在添加任务....")
		err = shell.CurrentUser.AddOfflineTask(identity, func() string {
			if tar, ok := c.Keys["o"]; ok {
				return tar
			}
			return shell.CurrentPath
		}())
		if err != nil {
			fmt.Println("添加失败:", err)
			return
		}
		fmt.Println("操作完成.")
	case "filter":
		if len(c.Nokeys) == 1 {
			fmt.Println("使用方法: offline filter <[-n 文件名] [-l 链接] [-s 状态] [-e 显示错误信息]>")
			fmt.Println("过滤离线任务列表 模糊搜索.")
			return
		}
	}
}
