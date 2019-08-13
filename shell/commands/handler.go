package commands

import (
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/shell"
	"github.com/Mrs4s/six-cli/six_cloud"
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
	shell.App.SetPrompt(shell.CurrentUser.Username + "@six-pan:" + models.ShortPath(shell.CurrentPath, 30) + "$ ")
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
