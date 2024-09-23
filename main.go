// Package main stands for default entrance
package main

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/horsing/ws/pkg"
	"github.com/horsing/ws/pkg/code"
	"github.com/horsing/ws/pkg/config"
	"github.com/horsing/ws/pkg/types"
	"github.com/horsing/ws/pkg/utils"
)

// Version semantic version string
var Version string = "dev"

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

type command struct {
	object  types.Application
	program string
	app     config.App
}

func main() {
	cfg := config.Get()

	if len(os.Args) <= 1 {
		usage(cfg)
		return
	}

	var commands []command
	var osargs []string

	for i, cmd := range os.Args[1:] {
		switch cmd {
		case "help", "-h", "--help", "?", "-?":
			usage(cfg)
			for _, sub := range os.Args[i+1:] {
				if _, ok := cfg.Programs[sub]; ok {
					fmt.Printf("  %+v\n", cfg.CommandHelp(sub))
				}
			}
			return
		case "-V", "--verbose":
			utils.Print("Enable verbose output")
		case "version":
			utils.Print(Version)
			return
		case "config":
			if vsc, ok := cfg.Programs["code"]; ok {
				code.New().Start(vsc.Program, config.Merge(cfg.Env, vsc.Env, insensitive()), []string{config.Configuration()}, vsc.Args...)
			} else {
				utils.Print("Program code is required but not configured.")
			}
			return
		default:
			if v, ok := cfg.Programs[cmd]; ok {
				switch cmd {
				case "code", "codium":
					commands = append(commands, command{code.New(), utils.Or(v.Program, cmd), v})
				default:
					commands = append(commands, command{pkg.GenericApplication{}, v.Program, v})
				}
			} else {
				osargs = append(osargs, cmd)
			}
		}
	}

	if len(commands) > 0 {
		for _, cmd := range commands {
			utils.Log("Starting %s with %v %v", cmd.program, osargs, cmd.app.Args)
			if err := cmd.object.Start(cmd.program, config.Merge(cfg.Env, cmd.app.Env, insensitive()), osargs, cmd.app.Args...); err != nil {
				utils.Error(err)
			}
		}
	}

	usage(cfg)
}
