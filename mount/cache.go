package mount

import (
	"errors"
	"fmt"
	"github.com/Mrs4s/six-cli/six_cloud"
	"io/ioutil"
	"math"
	"net/http"
)

var ErrDirUnsupported = errors.New("dir unsupported")

// Size of one chunk
var ChunkSize int64 = 1024 * 1024

// Max chunks cached on RAM
var ChunkMax int64 = 6

// Chunk preload count
var ChunkPreload int64

type FileCache struct {
	File   *six_cloud.SixFile
	Chunks map[int64][]byte

	queue chan CacheRequest
}

type CacheRequest struct {
	Request  *ReadRequest
	Callback chan ReadCallback
}

type ReadRequest struct {
	ChunkId     int64
	ChunkOffset int64
	Size        int64
}

type ReadCallback struct {
	Error   error
	Payload []byte
}

func NewCache(file *six_cloud.SixFile) (*FileCache, error) {
	if file.IsDir {
		return nil, ErrDirUnsupported
	}
	c := &FileCache{
		File:   file,
		Chunks: make(map[int64][]byte),
		queue:  make(chan CacheRequest),
	}
	go c.loop()
	return c, nil
}

func (c *FileCache) loop() {
	downloadQueue := make(chan CacheRequest)
	go func() {
		do := func(r CacheRequest) {
			if b, ok := c.Chunks[r.Request.ChunkId]; ok {
				if r.Callback != nil {
					begin := int64(math.Min(float64(r.Request.ChunkOffset), float64(len(b))))
					end := int64(math.Min(float64(r.Request.ChunkOffset+r.Request.Size), float64(len(b))))
					r.Callback <- ReadCallback{
						Payload: b[begin:end],
					}
					close(r.Callback)
				}
				return
			}
			chunkStart, chunkEnd := getChunk(r.Request.ChunkId)
			addr, err := c.File.GetDownloadAddress()
			if err != nil {
				r.callErr(err)
				return
			}
			b, err := download(addr, chunkStart, chunkEnd)
			if err != nil {
				r.callErr(err)
				return
			}
			if int64(len(c.Chunks)) > ChunkMax {
				var del int64 = -1
				for k := range c.Chunks {
					if int64(math.Abs(float64(k-r.Request.ChunkId))) > ChunkMax/2 {
						del = k
						break
					}
				}
				if del != -1 {
					delete(c.Chunks, del)
				}
			}
			c.Chunks[r.Request.ChunkId] = b
			if r.Callback != nil {
				begin := int64(math.Min(float64(r.Request.ChunkOffset), float64(len(b))))
				end := int64(math.Min(float64(r.Request.ChunkOffset+r.Request.Size), float64(len(b))))
				r.Callback <- ReadCallback{
					Payload: b[begin:end],
				}
				close(r.Callback)
			}
		}
		for r := range downloadQueue {
			go do(r)
			do(<-downloadQueue)
		}
	}()
	for req := range c.queue {
		if b, ok := c.Chunks[req.Request.ChunkId]; ok {
			if req.Callback != nil {
				begin := int64(math.Min(float64(req.Request.ChunkOffset), float64(len(b))))
				end := int64(math.Min(float64(req.Request.ChunkOffset+req.Request.Size), float64(len(b))))
				req.Callback <- ReadCallback{
					Payload: b[begin:end],
				}
				close(req.Callback)
			}
			continue
		}
		downloadQueue <- req
	}
	close(downloadQueue)
}

func (c *FileCache) Read(offset, size int64, callback chan ReadCallback) {
	chunkId := offset / ChunkSize
	fmt.Println("[", c.File.Name, "] load chunk", chunkId)
	req := &ReadRequest{
		ChunkId:     chunkId,
		ChunkOffset: offset % ChunkSize,
		Size:        size,
	}
	c.queue <- CacheRequest{
		Request:  req,
		Callback: callback,
	}
	for i := 1; i <= int(ChunkPreload); i++ {
		fmt.Println("[", c.File.Name, "] preload chunk", chunkId+int64(i))
		_, e := getChunk(chunkId + int64(i))
		if e <= c.File.Size {
			go func() {
				c.queue <- CacheRequest{
					Request: &ReadRequest{
						ChunkId: chunkId + int64(i),
					},
				}
			}()
		}
	}

}

func getChunk(chunkOffset int64) (chunkStart, chunkEnd int64) {
	chunkStart = chunkOffset * ChunkSize
	chunkEnd = chunkStart + ChunkSize
	return
}

func (req CacheRequest) callErr(err error) {
	if req.Callback != nil {
		req.Callback <- ReadCallback{
			Error: err,
		}
		close(req.Callback)
	}
}

func download(url string, beginOffset, endOffset int64) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", "Six-cli fuse")
	request.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", beginOffset, endOffset))
	rsp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
