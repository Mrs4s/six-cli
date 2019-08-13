package commands

import (
	"fmt"
	pl "github.com/Mrs4s/power-liner"
	"github.com/Mrs4s/six-cli/shell"
)

func init() {
	alias["Pwd"] = []string{}
	explains["Pwd"] = "显示当前工作目录"
}

func (CommandHandler) Pwd(c *pl.Context) {
	fmt.Println(shell.CurrentPath)
}
