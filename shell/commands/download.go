package commands

import (
	"fmt"
	httpDownloader "github.com/Mrs4s/go-http-downloader"
	pl "github.com/Mrs4s/power-liner"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/shell"
	"github.com/Mrs4s/six-cli/six_cloud"
	"github.com/cheggaaa/pb"
	"strconv"
	"strings"
	"time"
)

func init() {
	alias["Download"] = []string{"down"}
	explains["Download"] = "下载 文件夹/文件"
}

func (CommandHandler) Download(c *pl.Context) {
	if len(c.Nokeys) == 0 {
		fmt.Println("[H] 使用方法: down <文件/目录>")
		return
	}
	path := c.Nokeys[0]
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	targetPath := shell.CurrentPath
	if len(path) == 0 {
		fmt.Println("[H] 使用方法: down <文件/目录>")
		return
	}
	if strings.HasPrefix(path, "/") {
		targetPath = models.GetParentPath(path)
	}
	files, err := shell.CurrentUser.GetFilesByPath(targetPath)
	if err != nil {
		fmt.Println("[!] 错误:", err)
		return
	}
	var target *six_cloud.SixFile
	for _, file := range files {
		if strings.HasPrefix(file.Name, models.GetFileName(path)) {
			target = file
			continue
		}
	}
	if target == nil {
		fmt.Println("[!] 错误: 目标文件/目录不存在")
		return
	}
	var downloaders []*httpDownloader.DownloaderClient
	for key, file := range target.GetLocalTree(models.DefaultConf.DownloadPath) {
		fmt.Println("[+] 添加下载", models.ShortString(models.GetFileName(file.Path), 70))
		addr, err := file.GetDownloadAddress()
		if err != nil {
			fmt.Println("[!] 获取文件", file.Name, "的下载链接失败:", err)
			continue
		}
		info, err := httpDownloader.NewDownloaderInfo([]string{addr}, key, models.DefaultConf.DownloadBlockSize, int(models.DefaultConf.DownloadThread),
			map[string]string{"User-Agent": "Six-cli download engine"})
		client := httpDownloader.NewClient(info)
		p := file
		client.RefreshFunc = func() []string {
			addr, err := p.GetDownloadAddress()
			if err != nil {
				return []string{}
			}
			return []string{addr}
		}
		client.RefreshTime = 1000 * 60 * 3
		downloaders = append(downloaders, client)
	}
	ch := make(chan bool)
	defer close(ch)
	var bars []*pb.ProgressBar
	for _, task := range downloaders {
		bar := pb.New64(task.Info.ContentSize).Prefix(models.ShortString(models.GetFileName(task.Info.TargetFile), 20)).SetUnits(pb.U_BYTES)
		bar.ShowSpeed = true
		bars = append(bars, bar)
	}
	go func() {
		ticker := time.NewTicker(time.Second).C
		for range ticker {
			downloadingCount := 0
			waitingTask := -1
			for i, task := range downloaders {
				if task.Downloading {
					downloadingCount++
					bars[i].Set64(task.DownloadedSize)
				}
				if waitingTask == -1 && !task.Downloading && !task.Completed {
					waitingTask = i
				}
			}
			if downloadingCount < 1 && waitingTask != -1 {
				task := downloaders[waitingTask]
				fmt.Println()
				fmt.Println("[+] 即将开始下载任务 " + strconv.FormatInt(int64(waitingTask), 10))
				fmt.Println("[+] 文件名: " + models.GetFileName(task.Info.TargetFile))
				fmt.Println("[+] 下载路径: " + task.Info.TargetFile)
				fmt.Println("[+] 文件大小: " + models.ConvertSizeString(task.Info.ContentSize))
				fmt.Println()
				bars[waitingTask].Start()
				err := task.BeginDownload()
				if err != nil {
					bars[waitingTask].Finish()
					fmt.Println("[-] 文件", models.GetFileName(task.Info.TargetFile), "下载失败:", err)
					continue
				}
				task.OnCompleted(func() {
					bars[waitingTask].Finish()
					fmt.Println("[+] 文件", models.GetFileName(task.Info.TargetFile), "下载完成")
				})
				task.OnFailed(func(err error) {
					bars[waitingTask].Finish()
					fmt.Println("[-] 文件", models.GetFileName(task.Info.TargetFile), "下载失败:", err)
				})
			}
			if downloadingCount == 0 && waitingTask == -1 {
				ch <- true
				break
			}
		}
	}()
	<-ch
	time.Sleep(time.Second)
	fmt.Println("[+] 所有文件已下载完成.")
}

func (CommandHandler) DownloadCompleter(c *pl.Context) []string {
	return PathCompleter(c, true)
}
