package commands

import (
	"bytes"
	"fmt"
	pl "github.com/Mrs4s/power-liner"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/shell"
	"github.com/Mrs4s/six-cli/six_cloud"
	"github.com/jackpal/bencode-go"
	"strconv"
	"strings"
	"time"
)

func init() {
	alias["Preview"] = []string{"pw"}
}

func (CommandHandler) Preview(c *pl.Context) {
	if len(c.Nokeys) != 1 {
		fmt.Println("[H] 使用方法: pw <文件>")
		return
	}
	path := c.Nokeys[0]
	if !strings.HasPrefix(path, "/") {
		path = models.CombinePaths(shell.CurrentPath, path, "/")
	}
	file, err := shell.CurrentUser.GetFileByPath(path)
	if err != nil {
		fmt.Println("[!] 获取失败:", err)
		return
	}
	fmt.Println("[I] 文件名:", file.Name)
	fmt.Println("[I] 绝对路径:", file.Path)
	fmt.Println("[I] 文件大小:", file.Size, "("+models.ConvertSizeString(file.Size)+")")
	fmt.Println("[I] 文件夹:", file.IsDir)
	fmt.Println("[I] 创建时间:", time.Unix(file.CreateTime/1000, 0).In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05"))
	fmt.Println("[I] Mime:", file.Mime)
	if !file.IsDir {
		fmt.Println("[I] Hash:", file.ETag)
	}
	switch strings.ToLower(models.GetFileExtension(file.Name)) {
	case "txt", "ini", "inf":
		previewText(file)
	case "torrent":
		previewTorrent(file)
	}
}

func previewText(file *six_cloud.SixFile) {
	if file.Size < 1024*1024 {
		addr, err := file.GetDownloadAddress()
		if err == nil {
			text := shell.CurrentUser.Client.GetString(addr)
			if text != "" {
				fmt.Println()
				fmt.Println("[I] 预览: ")
				fmt.Println()
				lines := strings.Split(text, "\n")
				for i, line := range lines {
					fmt.Println("["+strconv.FormatInt(int64(i+1), 10)+"]:", line)
				}
			}
		}
	}
}

func previewTorrent(file *six_cloud.SixFile) {
	if file.Size < 1024*1024*5 {
		addr, err := file.GetDownloadAddress()
		if err == nil {
			b, err := shell.CurrentUser.Client.GetBytes(addr)
			if err == nil {
				fmt.Println()
				fmt.Println("[I] 预览: ")
				fmt.Println()
				dir, err := bencode.Decode(bytes.NewReader(b))
				if err != nil {
					fmt.Println("[!] 生成预览失败:", err)
					return
				}
				info, ok := dir.(map[string]interface{})["info"]
				if !ok {
					fmt.Println("[!] 生成预览失败: INFO字段不存在")
					return
				}
				fmt.Println("[T] Torrent名:", info.(map[string]interface{})["name"])
				if anno, ok := dir.(map[string]interface{})["announce"]; ok {
					fmt.Println("[T] Announce:", anno)
				}
				if commt, ok := dir.(map[string]interface{})["comment"]; ok {
					fmt.Println("[T] 评论:", commt)
				}
				if cb, ok := dir.(map[string]interface{})["created by"]; ok {
					fmt.Println("[T] 创建工具:", cb)
				}
				if sou, ok := info.(map[string]interface{})["source"]; ok {
					fmt.Println("[T] 来源:", sou)
				}
				if files, ok := info.(map[string]interface{})["files"]; ok {
					var size int64
					for _, sub := range files.([]interface{}) {
						if length, ok := sub.(map[string]interface{})["length"]; ok {
							size += length.(int64)
						}
					}
					fmt.Println("[T] 文件数量", len(files.([]interface{})), "个, 总大小", models.ConvertSizeString(size))
				} else if length, ok := info.(map[string]interface{})["length"]; ok {
					fmt.Println("[T] 文件数量 1 个, 总大小", models.ConvertSizeString(length.(int64)))
				}
			}
		}
	}
}

func (CommandHandler) PreviewCompleter(c *pl.Context) []string {
	return PathCompleter(c, true)
}
