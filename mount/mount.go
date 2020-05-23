// +build !windows

package mount

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"errors"
	fst "github.com/Mrs4s/six-cli/models/fs"
	"github.com/Mrs4s/six-cli/six_cloud"
	"os"
)

func Mount(user *six_cloud.SixUser, mountPoint string) error {
	if !fst.PathExists(mountPoint) {
		if err := os.MkdirAll(mountPoint, 0644); nil != err {
			return errors.New("could not create mount directory")
		}
	}
	options := []fuse.MountOption{
		fuse.AllowOther(),
		fuse.AllowNonEmptyMount(),
		fuse.ReadOnly(),
		fuse.FSName("six-pan"),
	}
	c, err := fuse.Mount(mountPoint, options...)
	if err != nil {
		return err
	}
	defer c.Close()
	filesys := &SixFileSystem{user: user}
	if err := fs.Serve(c, filesys); err != nil {
		return err
	}
	Unmount(mountPoint)
	return nil
}

func Unmount(mountPoint string) {
	_ = fuse.Unmount(mountPoint)
}
