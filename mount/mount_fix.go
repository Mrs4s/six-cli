// +build windows

package mount

import (
	"errors"
	"github.com/Mrs4s/six-cli/six_cloud"
	"runtime"
)

func Mount(user *six_cloud.SixUser, mountPoint string) error {
	return errors.New("mount unsupported " + runtime.GOOS)
}
