// +build !windows

package mount

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"context"
	"fmt"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/six_cloud"
	"os"
)

type SixFileSystem struct {
	user *six_cloud.SixUser
}

type FileObject struct {
	user  *six_cloud.SixUser
	file  *six_cloud.SixFile
	cache *FileCache
}

func (fs *SixFileSystem) Root() (fs.Node, error) {
	return &FileObject{
		user: fs.user,
		file: fs.user.GetRootFile(),
	}, nil
}

func (o *FileObject) Attr(ctx context.Context, attr *fuse.Attr) error {
	if o.file.IsDir {
		attr.Mode = os.ModeDir | 0755
		attr.Size = 0
	} else {
		attr.Mode = 0644
		attr.Size = uint64(o.file.Size)
	}
	attr.Blocks = (attr.Size + 511) / 512
	return nil
}

func (o *FileObject) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	if o.file.IsDir {
		return fuse.EIO
	}
	if o.cache == nil {
		o.cache, _ = NewCache(o.file)
	}
	callback := make(chan ReadCallback)
	o.cache.Read(req.Offset, int64(req.Size), callback)
	rsp := <-callback
	if rsp.Error != nil {
		fmt.Println("read error:", rsp.Error)
		return fuse.EIO
	}
	resp.Data = rsp.Payload
	return nil
}

func (o *FileObject) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	if !o.file.IsDir {
		return nil, fuse.ENOENT
	}
	var dirs []fuse.Dirent
	for _, sub := range o.file.GetChildren() {
		if sub.IsDir {
			dirs = append(dirs, fuse.Dirent{
				Type: fuse.DT_Dir,
				Name: sub.Name,
			})
			continue
		}
		dirs = append(dirs, fuse.Dirent{
			Type: fuse.DT_File,
			Name: sub.Name,
		})
	}
	return dirs, nil
}

func (o *FileObject) Lookup(ctx context.Context, name string) (fs.Node, error) {
	path := ""
	if o.file.Path == "/" {
		path = "/" + name
	} else {
		path = o.file.Path + "/" + name
	}
	file, err := o.user.GetFileByPath(path)
	if err != nil {
		return nil, fuse.ENOENT
	}
	return &FileObject{
		user: o.user,
		file: file,
	}, nil
}

func (o *FileObject) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	err := o.user.DeleteFile(o.file.Path)
	if err != nil {
		return fuse.EIO
	}
	return nil
}

func (o *FileObject) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	path := models.CombinePaths(o.file.Path, req.Name, "/")
	err := o.user.CreateDirectory(path)
	if err != nil {
		return nil, fuse.EIO
	}
	return &FileObject{
		user: o.user,
		file: &six_cloud.SixFile{
			Path:  path,
			Name:  req.Name,
			IsDir: true,
		},
	}, nil
}

/*
func (o *FileObject) Rename(ctx context.Context, req *fuse.RenameRequest, newDir fs.Node) error {
}
*/
