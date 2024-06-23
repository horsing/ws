package pkg

import (
	"os/exec"
)

type GenericApplication struct {
}

func (g GenericApplication) Start(program string, env []string, osargs []string, args ...string) bool {
	varg := append([]string{}, args...)
	varg = append(varg, osargs...)
	c := exec.Command(program, varg...)
	c.Env = env
	return c.Start() == nil
}
