package code

import (
	"github.com/horsing/coder/pkg"
	"runtime"
)

func New() pkg.Application {
	if runtime.GOOS == "windows" {
		return windows{}
	}
	return linux{}
}
