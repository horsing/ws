package code

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

type windows struct {
}

const folder = "Microsoft VS Code"

var programs = []string{
	path.Join(os.Getenv("ProgramFiles"), folder),
	path.Join(os.Getenv("ProgramFiles(x86)"), folder),
	path.Join(os.Getenv("LOCALAPPDATA"), "Programs", folder),
}

const p = "Code.exe"

func (w windows) Start(program string, env []string, osargs []string, args ...string) bool {
	font := args[0]

	if len(program) == 0 {
		for _, d := range programs {
			if _, err := os.Stat(path.Join(d, p)); err == nil {
				program = d
				break
			}
		}
	}

	location := program[:strings.LastIndex(program, string(os.PathSeparator))]
	if len(location) == 0 {
		return false
	}

	css := path.Join(location, "resources/app/out/vs/workbench/workbench.desktop.main.css")
	js := path.Join(location, "resources/app/out/vs/workbench/workbench.desktop.main.js")
	regcss := regexp.MustCompile("font-family:([^,]+, *)?Segoe WPC,")
	regjs := regexp.MustCompile("font-family: ([^,]+, *)?\"Segoe WPC\",")

	if buf, err := os.ReadFile(css); err == nil {
		body := string(buf)
		replaced := regcss.ReplaceAllString(body, fmt.Sprintf("font-family:%s, Segoe WPC,", font))
		if err := os.WriteFile(css, []byte(replaced), fs.ModeType); err != nil {
			return false
		}
	}

	if buf, err := os.ReadFile(js); err == nil {
		body := string(buf)
		replaced := regjs.ReplaceAllString(body, fmt.Sprintf("font-family: \"%s\", \"Segoe WPC\",", font))
		if err := os.WriteFile(js, []byte(replaced), fs.ModeType); err != nil {
			return false
		}
	}

	c := exec.Command(program, osargs...)
	c.Env = env
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	if err := c.Start(); err == nil {
		return true
	} else {
		fmt.Println(err)
		return false
	}
}
