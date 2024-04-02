package code

import "runtime"

type Code interface {
	Replace(font string, dirs ...string) bool
}

func New() Code {
	if runtime.GOOS == "windows" {
		return windows{}
	}
	return linux{}
}