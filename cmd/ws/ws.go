package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/horsing/ws/pkg"
	"github.com/horsing/ws/pkg/code"
	"github.com/horsing/ws/pkg/config"
)

func usage(c *config.Config) {
	_, cmd := path.Split(os.Args[0])
	if runtime.GOOS == "windows" {
		if i := strings.LastIndex(cmd, string(os.PathSeparator)); i >= 0 {
			cmd = cmd[i+1:]
		}
	}
	fmt.Printf("Usage: %s {help|version|config|%s}\n", cmd, c.AvailableCommands("|"))
}

func insensitive() bool {
	return runtime.GOOS == "windows"
}

func main() {
	cfg := config.Get()

	if len(os.Args) <= 1 {
		usage(cfg)
		return
	}

	switch os.Args[1] {
	case "help":
		usage(cfg)
		if len(os.Args) > 2 {
			v := os.Args[2]
			if _, ok := cfg.Programs[v]; ok {
				fmt.Printf("  %v", cfg.CommandHelp(v))
			}
		}
	case "version":
		fmt.Println("dev")
	case "config":
		if vsc, ok := cfg.Programs["code"]; ok {
			code.New().Start(vsc.Program, config.Merge(cfg.Env, vsc.Env, insensitive()), []string{config.Configuration()}, vsc.Args...)
		} else {
			panic("Program code is required but not configured.")
		}
	default:
		if v, ok := cfg.Programs[os.Args[1]]; ok {
			switch os.Args[1] {
			case "code":
				code.New().Start(v.Program, config.Merge(cfg.Env, v.Env, insensitive()), os.Args[2:], v.Args...)
			default:
				if err := (pkg.GenericApplication{}.Start(v.Program, config.Merge(cfg.Env, v.Env, insensitive()), os.Args[1:], v.Args...)); err != nil {
					panic(err)
				}
			}
		} else {
			usage(cfg)
		}
	}
}
