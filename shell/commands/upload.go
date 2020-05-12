package commands

import (
	"fmt"
	sixcloudUploader "github.com/Mrs4s/go-six-cloud-upload-sdk"
	pl "github.com/Mrs4s/power-liner"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/models/fs"
	"github.com/Mrs4s/six-cli/shell"
	"os"
	"time"
)

func init() {
	alias["Upload"] = []string{"up"}
	explains["Upload"] = "上传文件 / 文件夹"
}

func (CommandHandler) Upload(c *pl.Context) {
	if len(c.Nokeys) == 0 {
		fmt.Println("[?] 使用方法: upload <本地文件> [-o 远程目录]")
		fmt.Println("[?] 上传本地文件或文件夹到远程目录, 默认工作目录")
		return
	}
	local := c.Nokeys[0]
	remote := shell.CurrentPath
	if r, ok := c.Keys["o"]; ok {
		remote = r
	}
	if !models.PathExists(local) {
		fmt.Println("[!] 本地文件不存在.")
		return
	}
	if _, err := shell.CurrentUser.GetFileByPath(remote); err != nil {
		if remote != "/" {
			fmt.Println("[!] 远程目录不存在.")
			return
		}
	}
	tree := shell.CurrentUser.CreateUploadTree(remote, c.Nokeys)
	fmt.Println("[+] 正在准备上传", len(tree), "个文件.")
	fmt.Println()
	for k, v := range tree {
		fmt.Println("[+] 正在读取文件", fs.GetFileName(v))
		etag, err := fs.ComputeFileEtag(v)
		if err != nil {
			fmt.Println("[!] 读取出错:", err)
			fmt.Println()
			continue
		}
		token := shell.CurrentUser.CreateUploadToken(fs.GetParentPath(k), fs.GetFileName(k), etag)
		if token.Cached {
			fmt.Println("[+] 文件", fs.GetFileName(v), "秒传完成.")
			fmt.Println()
			continue
		}
		if token.UploadToken == "" {
			fmt.Println("[!] 获取上传Token失败.")
			fmt.Println()
			return
		}
		info, _ := sixcloudUploader.CreateUploadTask(token.UploadToken, token.UploadUrl, v)
		client := sixcloudUploader.NewClient(info, 4)
		uploadedBlocks := func() (num int) {
			for _, b := range client.Info.Blocks {
				if b.Uploaded {
					num++
				}
			}
			return
		}
		client.BeginUpload()
		// TODO: Need fix.
		client.OnUploaded = func(client *sixcloudUploader.UploadClient) {
			client.Status = sixcloudUploader.Completed
		}
		client.OnUploadFailed = func(client *sixcloudUploader.UploadClient) {
			client.Status = sixcloudUploader.Failed
		}
		for client.Status == sixcloudUploader.Uploading {
			time.Sleep(time.Second)
			_, _ = fmt.Fprintf(os.Stdout, "\r ↑ %d / %d (分块)", uploadedBlocks(), len(info.Blocks))
		}
		fmt.Println()
		fmt.Println()
		fmt.Println("[+] 文件", fs.GetFileName(v), "处理完成.")
	}
	fmt.Println("[+] 所有文件已处理完成.")
}
