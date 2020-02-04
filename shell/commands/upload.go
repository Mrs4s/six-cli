package commands

import (
	"fmt"
	pl "github.com/Mrs4s/power-liner"
)

func init() {
	alias["Upload"] = []string{"up"}
	explains["Upload"] = "上传文件 / 文件夹"
}

func (CommandHandler) Upload(c *pl.Context) {
	if len(c.Nokeys) == 0 {
		fmt.Println("[H] 使用方法: upload <本地文件> [-o 远程目录]")
		fmt.Println("上传本地文件或文件夹到远程目录, 默认工作目录")
		return
	}
}
