package code

import (
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strings"
)

type windows struct{}

const folder = "Microsoft VS Code"

var programs = []string{path.Join(os.Getenv("ProgramFiles"), folder),
	path.Join(os.Getenv("ProgramFiles(x86)"), folder),
	path.Join(os.Getenv("LOCALAPPDATA"), "Programs", folder),
}

func (w windows) Replace(font string, dirs ...string) bool {
	locations := append(dirs, programs...)
	var location string
	for _, l := range locations {
		if _, err := os.Stat(l); err == nil {
			location = l
			break
		}
	}
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

	_ = exec.Command(path.Join(location, "Code.exe")).Start()
	return true
}