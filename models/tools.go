package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func ToMd5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
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

func CombinePaths(path1, path2, sep string) string {
	if len(path2) == 0 {
		return path1
	}
	if len(path1) == 0 {
		return path2
	}
	char := path1[len(path1)-1:]
	if sep == "" {
		sep = string(os.PathSeparator)
	}
	if char != "\\" && char != "/" && char != ":" {
		return path1 + sep + path2
	}
	return path1 + path2
}

func ConvertSizeString(size int64) string {
	switch {
	case size <= 0:
		return "0B"
	case size <= 1024: // B
		return fmt.Sprintf("%dB", size)
	case size < 1024*1024: // KB
		return fmt.Sprintf("%.2fKB", float64(size)/float64(1024))
	case size < 1024*1024*1024: // MB
		return fmt.Sprintf("%.2fMB", float64(size)/float64(1024*1024))
	case size < 1024*1024*1024*1024: //GB
		return fmt.Sprintf("%.2fGB", float64(size)/float64(1024*1024*1024))
	default:
		return fmt.Sprintf("%.2fTB", float64(size)/float64(1024*1024*1024*1024))
	}
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
