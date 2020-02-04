package commands

import (
	pl "github.com/Mrs4s/power-liner"
)

func init() {
	alias["Copy"] = []string{"cp"}
}

func (CommandHandler) Copy(c *pl.Context) {
	// TODO: bug fix
	/*
		args := models.FilterStrings(c.Nokeys, func(s string) bool { return s != "" })
		if len(args) != 2 {
			fmt.Println("使用方法: cp <源文件> <目标文件> [-f]")
			return
		}
		if _, err := shell.CurrentUser.GetFileByPath(args[1]); err == nil {
			if _, ok := c.Keys["f"]; !ok {
				if i, _ := shell.App.ReadLine("目标文件已存在，是否覆盖? (y/n)"); i != "y" {
					return
				}
			}
		}
		err := shell.CurrentUser.CopyFile(args[0], args[1])
		if err != nil {
			fmt.Println("操作出错:", err)
			return
		}
		fmt.Println("操作完成.")
	*/
}

func (CommandHandler) CopyCompleter(c *pl.Context) []string {
	return PathCompleter(c, true)
}
