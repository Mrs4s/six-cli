package commands

import (
	"fmt"
	pl "github.com/Mrs4s/power-liner"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/mount"
	"github.com/Mrs4s/six-cli/shell"
)

func init() {
	alias["Mount"] = []string{}
	explains["Mount"] = "挂载目录"
}

func (CommandHandler) Mount(c *pl.Context) {
	args := models.FilterStrings(c.Nokeys, func(s string) bool { return s != "" })
	if len(args) != 1 {
		fmt.Println("[?] 使用方法: Mount <挂载点>")
		return
	}
	p := args[0]
	mount.ChunkPreload = 4
	if err := mount.Mount(shell.CurrentUser, p); err != nil {
		fmt.Println("[!] 错误:", err)
	}
}
