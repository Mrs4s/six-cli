package commands

import (
	"fmt"
	pl "github.com/Mrs4s/power-liner"
)

func init() {

}

func (CommandHandler) Switch(c *pl.Context) {
	if len(c.RawArgs) == 1 {
		fmt.Println("用户列表:")
		printUserList()
		return
	}
}
