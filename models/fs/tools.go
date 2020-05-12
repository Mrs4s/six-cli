package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func GetFileName(path string) string {
	length := len(path)
	index := length - 1
	for index > 0 {
		char := path[index-1 : index]
		if char == "\\" || char == "/" || char == ":" {
			return path[index:length]
		}
		index--
	}
	return path
}

func IsDir(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.IsDir()
}

func GetDirEntities(path string) []string {
	entities, err := ioutil.ReadDir(path)
	if err != nil {
		return []string{}
	}
	var res []string
	for _, file := range entities {
		res = append(res, filepath.Join(path, file.Name()))
	}
	return res
}

func GetParentPath(path string) string {
	list := strings.Split(path, "/")
	var tmp []string
	for i := 0; i < len(list)-1; i++ {
		tmp = append(tmp, list[i])
	}
	parentPath := strings.Join(tmp, "/")
	if parentPath == "" {
		return "/"
	}
	return parentPath
}
