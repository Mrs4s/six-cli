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
	alias["List"] = []string{"ls"}
	explains["List"] = "显示 当前目录/目标目录 下的所有文件"
}

func (CommandHandler) List(c *pl.Context) {
	var (
		//blue       = color.New(color.FgBlue).SprintFunc()
		//green      = color.New(color.FgGreen).SprintFunc()
		files, err = shell.CurrentUser.GetFilesByPath(shell.CurrentPath)
	)
	if len(c.Nokeys) >= 1 {
		tar := strings.Join(c.Nokeys, " ")
		if len(tar) == 0 || tar[0:1] != "/" {
			fmt.Println("[-] 非法操作")
			return
		}
		files, err = shell.CurrentUser.GetFilesByPath(tar)
	}
	if err != nil {
		fmt.Println("[-] 获取文件失败: " + err.Error())
		return
	}
	if len(files) == 0 {
		fmt.Println("[-] 当前目录下无任何文件")
		return
	}
	printColor := func(files []*six_cloud.SixFile) (strs []string) {
		// 暂不支持染色
		for _, file := range files {
			if file.IsDir {
				//strs = append(strs, blue(file.Name))
				strs = append(strs, file.Name)
				continue
			}
			//strs = append(strs, green(file.Name))
			strs = append(strs, file.Name)
		}
		return
	}
	filteredFiles := files
	if _, ok := c.Keys["R"]; ok {
		fmt.Println(".:")
		shell.App.PrintColumns(printColor(filteredFiles), 2)
		for _, dir := range filterDirs(files) {
			subFiles, err := shell.CurrentUser.GetFilesByPath(models.CombinePaths(shell.CurrentPath, dir, "/"))
			if err == nil {
				fmt.Println("\n./" + dir + ":")
				shell.App.PrintColumns(printColor(append(subFiles)), 2)
			}
		}
		return
	}
	if _, ok := c.Keys["d"]; ok {
		var dirs []*six_cloud.SixFile
		for _, file := range files {
			if file.IsDir {
				dirs = append(dirs, file)
			}

		}
		filteredFiles = dirs
	}

	if _, ok := c.Keys["a"]; ok {
		table := [][]string{{"序号", "创建时间", "文件大小", "文件名"}}
		for i, file := range filteredFiles {
			table = append(table, []string{
				strconv.FormatInt(int64(i), 10),
				time.Unix(file.CreateTime/1000, 0).In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05"),
				models.ConvertSizeString(file.Size),
				file.Name,
			})
		}
		shell.App.PrintTables(table, 2)
		return
	}
	shell.App.PrintColumns(printColor(filteredFiles), 2)
}
