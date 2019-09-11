package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/mattn/go-runewidth"
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

func ShortString(str string, length int) string {
	sb := strings.Builder{}
	var num int
	for _, s := range str {
		num += runewidth.RuneWidth(s)
		if num > length {
			sb.WriteString("...")
			break
		}
		sb.WriteRune(s)
	}
	return sb.String()
}

func ShortPath(path string, length int) string {
	sub := FilterStrings(strings.Split(path, "/"), func(s string) bool { return s != "" })
	if len(sub) == 0 {
		return path
	}
	return "/" + strings.Join(SelectStrings(sub, func(s string) string { return ShortString(s, length) }), "/")
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func GetFileExtension(file string) string {
	sp := strings.Split(file, ".")
	return sp[len(sp)-1]
}

func ShellMatch(str, p string) bool {
	var (
		j, i, star, last int
		rs               = []rune(str)
		rp               = []rune(p)
	)

	for i < len(rs) {
		if j < len(rp) && (rs[i] == rp[j] || rp[j] == '?') {
			i++
			j++
		} else if j < len(rp) && rp[j] == '*' {
			j++
			last = i
			star = j
		} else if star != 0 {
			last++
			i = last
			j = star
		} else {
			return false
		}
	}
	for j < len(rp) && rp[j] == '*' {
		j++
	}
	return j == len(rp)
}

func SelectStrings(arr []string, selector func(string) string) []string {
	var res []string
	for _, str := range arr {
		res = append(res, selector(str))
	}
	return res
}

func FilterStrings(arr []string, filter func(string) bool) []string {
	var res []string
	for _, str := range arr {
		if filter(str) {
			res = append(res, str)
		}
	}
	return res
}
