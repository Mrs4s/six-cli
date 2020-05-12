package fs

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"os"
	"strconv"
)

func ComputeFileEtag(filename string) (etag string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return
	}

	fsize := fi.Size()
	innerBlockCount := blockCount(fsize)
	var tag []byte

	if innerBlockCount <= 1 { // file size <= 4M
		tag, err = computeSha1([]byte{0x16}, f)
		if err != nil {
			return
		}
	} else { // file size > 4M
		var allBlocksSha1 []byte

		for i := 0; i < innerBlockCount; i++ {
			body := io.LimitReader(f, csBlockSize)
			allBlocksSha1, err = computeSha1(allBlocksSha1, body)
			if err != nil {
				return
			}
		}

		tag, _ = computeSha1([]byte{0x96}, bytes.NewReader(allBlocksSha1))
	}

	etag = base64.URLEncoding.EncodeToString(tag)
	etag += strconv.FormatInt(fsize, 36)
	return
}

const (
	csBlockBits = 22               // 2 ^ 22 = 4M
	csBlockSize = 1 << csBlockBits // 4M
)

func blockCount(size int64) int {
	return int((size + (csBlockSize - 1)) >> csBlockBits)
}

func computeSha1(b []byte, r io.Reader) ([]byte, error) {
	h := sha1.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return nil, err
	}
	return h.Sum(b), nil
}
