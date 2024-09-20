package code

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/horsing/ws/pkg/utils"
)

type windows struct {
}

var folders = []string{"Microsoft VS Code", "VSCodium"}

var programs = []string{
	path.Join(os.Getenv("ProgramFiles")),
	path.Join(os.Getenv("ProgramFiles(x86)")),
	path.Join(os.Getenv("LOCALAPPDATA"), "Programs"),
	path.Join(os.Getenv("APPDATA"), "Apps"), // for VSCodium
}

func (w windows) Start(program string, env []string, osargs []string, args ...string) bool {
	p := "Code.exe"
	if program == "codium" {
		p = "VSCodium.exe"
	}

	font := args[0]
	location := ""

	if _, e := os.Stat(program); len(program) == 0 || e != nil {
		for _, d := range programs {
			for _, folder := range folders {
				ep := path.Join(d, folder, p)
				if _, err := os.Stat(ep); err == nil {
					location = path.Join(d, folder)
					program = ep
					break
				}
			}
		}
	} else {
		location = program[:strings.LastIndex(program, string(os.PathSeparator))]
	}

	if len(location) == 0 {
		return false
	}

	utils.Print("Starting %s in %s", p, location)

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
		utils.Print("Started [%s] with environment: %q", program, env)
		return true
	} else {
		utils.Error(err)
		return false
	}
}
