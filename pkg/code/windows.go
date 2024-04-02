package code

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strings"
)

type windows struct {
}

const folder = "Microsoft VS Code"

var programs = []string{path.Join(os.Getenv("ProgramFiles"), folder),
	path.Join(os.Getenv("ProgramFiles(x86)"), folder),
	path.Join(os.Getenv("LOCALAPPDATA"), "Programs", folder),
}

const p = "Code.exe"

func (w windows) Start(program string, env []string, args ...string) bool {
	font := args[0]
	location := program[:strings.LastIndex(program, string(os.PathSeparator))]
	if location == "" {
		return false
	}

	files := []string{
		"resources/app/out/vs/workbench/workbench.desktop.main.css",
		"resources/app/out/vs/workbench/workbench.desktop.main.js",
	}
	for _, file := range files {
		if buf, err := os.ReadFile(path.Join(location, file)); err == nil {
			body := string(buf)
			replaced := strings.ReplaceAll(body, "Segoe WPC", font)
			if err := os.WriteFile(path.Join(location, file), []byte(replaced), fs.ModeType); err != nil {
				return false
			}
		}
	}

	c := exec.Command(program)
	c.Env = env
	if err := c.Start(); err == nil {
		return true
	} else {
		fmt.Println(err)
		return false
	}
}
