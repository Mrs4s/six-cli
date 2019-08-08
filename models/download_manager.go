package models

import (
	"fmt"
	httpDownloader "github.com/Mrs4s/go-http-downloader"
	"time"
)

var (
	DownloaderList []*httpDownloader.DownloaderClient
)

func init() {
	go func() {
		ticker := time.NewTicker(time.Second).C
		for range ticker {
			if len(DownloaderList) > 0 {
				var downloadingCount int
				waitingTask := -1
				for i, task := range DownloaderList {
					if task.Downloading {
						downloadingCount++
					}
					if waitingTask == -1 && !task.Downloading {
						waitingTask = i
					}
				}
				if downloadingCount == 0 && waitingTask != -1 {
					task := DownloaderList[waitingTask]
					err := task.BeginDownload()
					if err != nil {
						fmt.Println("download error", err)
						DownloaderList = append(DownloaderList[:waitingTask], DownloaderList[waitingTask+1:]...)
					}
				}
			}
		}
	}()
}
