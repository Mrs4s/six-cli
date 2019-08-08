package six_cloud

import (
	"github.com/Mrs4s/six-cli/models"
	"os"
)

func (file *SixFile) GetDownloadAddress() (string, error) {
	return file.owner.GetDownloadAddressByPath(file.Path)
}

func (file *SixFile) GetChildren() []*SixFile {
	res, _ := file.owner.GetFilesByPath(file.Path)
	return res
}

func (file *SixFile) GetLocalTree(localPath string) map[string]*SixFile {
	if !models.PathExists(localPath) {
		_ = os.MkdirAll(localPath, os.ModePerm)
	}
	result := make(map[string]*SixFile)
	if !file.IsDir {
		result[models.CombinePaths(localPath, file.Name, "")] = file
		return result
	}
	for _, childrenFile := range file.GetChildren() {
		for path, node := range childrenFile.GetLocalTree(models.CombinePaths(localPath, file.Name, "")) {
			result[path] = node
		}
	}
	return result
}
