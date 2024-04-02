package main

import (
	"github.com/horsing/coder/pkg/code"
	"os"
)

func main() {
	code.New().Replace("JetBrains Mono", os.Args[1:]...)
}