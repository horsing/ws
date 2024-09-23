package types

type Application interface {
	Start(program string, env []string, osargs []string, args ...string) error
}
