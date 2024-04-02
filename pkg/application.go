package pkg

import (
	"os"
	"os/exec"
	"strings"
)

type App struct {
	Program string   `json:"program"`
	Args    []string `json:"args"`
	Env     []string `json:"env"`
}

type Config struct {
	Env      []string       `json:"env"`
	Programs map[string]App `json:"programs"`
}

type Application interface {
	Start(program string, env []string, args ...string) bool
}

type GenericApplication struct {
}

func (g GenericApplication) Start(program string, env []string, args ...string) bool {
	c := exec.Command(program)
	c.Env = env
	return c.Start() == nil
}

func Merge(lhs, rhs []string, insensitive bool) []string {
	merged := lhs
	for _, e := range rhs {
		if i := strings.Index(e, "="); i > 0 {
			matched := false
			for j, l := range merged {
				if len(l) > i && strings.EqualFold(l[:i+1], e[:i+1]) {
					// ...= matched
					merged[j] = l + string(os.PathListSeparator) + e[i+1:]
					matched = true
					break
				}
			}
			if !matched {
				merged = append(merged, e)
			}
		}
	}
	if len(merged) <= 0 {
		return lhs
	}
	return append(merged)
}
