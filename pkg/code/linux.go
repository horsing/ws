package code

type linux struct{}

func (l linux) Start(program string, env []string, osargs []string, args ...string) bool {
	panic("implement me")
}
