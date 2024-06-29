package code

import (
	"runtime"

	"github.com/horsing/ws/pkg/types"
)

func New() types.Application {
	if runtime.GOOS == "windows" {
		return windows{}
	}
	return linux{}
}
