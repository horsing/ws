package main

import (
	"fmt"
	"os"

	"github.com/horsing/ws/pkg"
	"github.com/horsing/ws/pkg/code"
	"github.com/horsing/ws/pkg/config"
)

func usage() {
	fmt.Printf(`Usage: %s <cmd>`, os.Args[0])
}

func main() {
	cfg := config.Get()

	if len(os.Args) <= 1 {
		usage()
		return
	}

	switch os.Args[1] {
	case "help":
		usage()
	case "version":
		fmt.Println("dev")
	case "config":
		if vsc, ok := cfg.Programs["code"]; ok {
			code.New().Start(vsc.Program, append(cfg.Env, vsc.Env...), []string{config.Configuration()}, vsc.Args...)
		} else {
			panic("Program code is required but not configured.")
		}
	default:
		if v, ok := cfg.Programs[os.Args[1]]; ok {
			switch os.Args[1] {
			case "code":
				code.New().Start(v.Program, append(cfg.Env, v.Env...), os.Args[2:], v.Args...)
			default:
				pkg.GenericApplication{}.Start(v.Program, append(cfg.Env, v.Env...), os.Args[1:], v.Args...)
			}
		} else {
			usage()
		}
	}
}
