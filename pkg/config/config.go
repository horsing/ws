package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

const workspace = ".workspace"

type App struct {
	Program string   `json:"program"`
	Args    []string `json:"args"`
	Env     []string `json:"env"`
}

type Config struct {
	Env      []string       `json:"env"`
	Programs map[string]App `json:"programs"`
}

func (c Config) AvailableCommands(sep string) string {
	keys := make([]string, len(c.Programs))
	i := 0
	for k := range c.Programs {
		keys[i] = k
		i += 1
	}
	return strings.Join(keys, sep)
}

func (c Config) CommandHelp(cmd string) []any {
	if v, ok := c.Programs[cmd]; ok {
		return []any{cmd, v}
	}
	return []any{}
}

func Configuration() string {
	if runtime.GOOS == "windows" {
		return path.Join(os.Getenv("USERPROFILE"), workspace)
	}
	return path.Join(os.Getenv("HOME"), workspace)
}

func Get() *Config {
	cfg := Config{
		Env:      os.Environ(),
		Programs: map[string]App{},
	}
	workspace := Configuration()
	if s, err := os.Stat(workspace); err == nil {
		if !s.IsDir() {
			if buf, err := os.ReadFile(workspace); err == nil {
				tmpCommands := Config{Env: []string{}, Programs: map[string]App{}}
				if err := json.Unmarshal(buf, &tmpCommands); err == nil {
					for k, v := range tmpCommands.Programs {
						cfg.Programs[k] = v
					}
					// windows insensitive
					cfg.Env = Merge(cfg.Env, tmpCommands.Env, runtime.GOOS == "windows")
				} else {
					fmt.Println(err)
				}
			}
		}
	}
	return &cfg
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
	return merged
}
