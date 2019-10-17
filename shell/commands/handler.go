package commands

import (
	"fmt"
	pl "github.com/Mrs4s/power-liner"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/shell"
	"github.com/Mrs4s/six-cli/six_cloud"
	"strings"
)

var (
	alias    = make(map[string][]string)
	explains = make(map[string]string)
)

type CommandHandler struct {
}

func (CommandHandler) Alias() map[string][]string {
	return alias
}

func (CommandHandler) Explains() map[string]string {
	return explains
}

func refreshPrompt() {
	if models.GetParentPath(shell.CurrentPath) != "/" {
		shell.App.SetPrompt(shell.CurrentUser.Username + "@six-pan:" + models.ShortPath(models.GetFileName(shell.CurrentPath), 30) + "$ ")
		return
	}
	shell.App.SetPrompt(shell.CurrentUser.Username + "@six-pan:" + models.ShortPath(shell.CurrentPath, 30) + "$ ")
}

func printUserList() {
	for i, user := range shell.SavedUsers {
		if user.Identity == shell.CurrentUser.Identity {
			fmt.Println(i+1, "->", user.Username)
			continue
		}
		fmt.Println(user.Username)
	}
}

func PathCompleter(c *pl.Context, f bool) []string {
	if len(c.Nokeys) > 1 {
		return []string{}
	}
	if len(strings.Split(c.Nokeys[0], "/")) <= 1 {
		fun := func(s string) string {
			if strings.Contains(s, " ") {
				return "\"" + s + "\""
			}
			if strings.Contains(s, ".") {
				return s
			}
			return s + "/"
		}
		if f {
			return models.SelectStrings(append(filterCurrentDirs(), filterCurrentFiles()...), fun)
		}
		return models.SelectStrings(filterCurrentDirs(), fun)
	}
	newPath := shell.CurrentPath + "/" + c.Nokeys[0]
	if shell.CurrentPath == "/" {
		newPath = "/" + c.Nokeys[0]
	}
	files, err := shell.CurrentUser.GetFilesByPath(models.GetParentPath(newPath))
	if err != nil {
		return []string{}
	}
	filter := func(s string) bool {
		if strings.HasSuffix(newPath, "/") {
			return true
		}
		return strings.HasPrefix(s, models.GetFileName(newPath))
	}
	selector := func(s string) string {
		com := models.CombinePaths(models.GetParentPath(newPath), s, "/")
		if strings.Contains(com, " ") {
			return "\"" + com[1:] + "\""
		}
		if strings.Contains(com[1:], ".") {
			return com[1:]
		}
		return com[1:] + "/"
	}
	fmt.Println(f)
	if f {
		return models.SelectStrings(models.FilterStrings(append(filterDirs(files), filterFiles(files)...), filter), selector)
	}
	return models.SelectStrings(models.FilterStrings(filterDirs(files), filter), selector)
}

func filterCurrentDirs() []string {
	files, err := shell.CurrentUser.GetFilesByPath(shell.CurrentPath)
	if err != nil {
		return []string{}
	}
	return filterDirs(files)
}

func filterDirs(files []*six_cloud.SixFile) (res []string) {
	for _, file := range files {
		if file.IsDir {
			res = append(res, file.Name)
		}
	}
	return
}

func filterCurrentFiles() []string {
	files, err := shell.CurrentUser.GetFilesByPath(shell.CurrentPath)
	if err != nil {
		return []string{}
	}
	return filterFiles(files)
}

func filterFiles(files []*six_cloud.SixFile) (res []string) {
	for _, file := range files {
		if !file.IsDir {
			res = append(res, file.Name)
		}
	}
	return
}

func filterFileArray(arr []*six_cloud.SixFile, filter func(*six_cloud.SixFile) bool) []*six_cloud.SixFile {
	var res []*six_cloud.SixFile
	for _, file := range arr {
		if filter(file) {
			res = append(res, file)
		}
	}
	return res
}
