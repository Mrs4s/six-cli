package commands

import pl "github.com/Mrs4s/power-liner"

func init() {
	alias["Config"] = []string{"cof"}
	explains["Config"] = "修改配置文件"
}

func (CommandHandler) Config(c *pl.Context) {

}
