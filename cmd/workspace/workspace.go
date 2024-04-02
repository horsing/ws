package main

import (
	"encoding/json"
	"fmt"
	"github.com/horsing/coder/pkg"
	"github.com/horsing/coder/pkg/code"
	"os"
	"path"
	"runtime"
)

func usage() {
	fmt.Printf(`Usage: %s <cmd>`, os.Args[0])
}

func main() {
	cfg := pkg.Config{
		Env:      os.Environ(),
		Programs: map[string]pkg.App{},
	}
	for _, dir := range []string{os.Getenv("HOME"), os.Getenv("USERPROFILE")} {
		if s, err := os.Stat(dir); err == nil {
			if s.IsDir() {
				workspace := path.Join(dir, ".workspace")
				if s, err := os.Stat(workspace); err == nil {
					if !s.IsDir() {
						if buf, err := os.ReadFile(workspace); err == nil {
							tmpCommands := pkg.Config{Env: []string{}, Programs: map[string]pkg.App{}}
							if err := json.Unmarshal(buf, &tmpCommands); err == nil {
								for k, v := range tmpCommands.Programs {
									cfg.Programs[k] = v
								}
								// windows insensitive
								cfg.Env = pkg.Merge(cfg.Env, tmpCommands.Env, runtime.GOOS == "windows")
							} else {
								fmt.Println(err)
							}
						}
					}
				}
			}
		}
	}

	if len(os.Args) <= 1 {
		usage()
		return
	}

	switch os.Args[1] {
	case "help":
		usage()
	case "version":
		fmt.Println("dev")
	default:
		if v, ok := cfg.Programs[os.Args[1]]; ok {
			switch os.Args[1] {
			case "code":
				code.New().Start(v.Program, append(cfg.Env, v.Env...), v.Args...)
			default:
				pkg.GenericApplication{}.Start(v.Program, append(cfg.Env, v.Env...), v.Args...)
			}
		} else {
			usage()
		}
	}
}
