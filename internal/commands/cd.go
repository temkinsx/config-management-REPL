package commands

import (
	"github.com/temkinsx/config-management-REPL/internal/commands/model"
)

type Cd struct{}

func (c *Cd) Name() string {
	return "cd"
}

func (c *Cd) Run(args []string, env *model.Env) (string, error) {
	if len(args) < 1 {
		env.FS.Cwd = env.FS.Root
		return "", nil
	}

	err := env.FS.Cd(args[0])
	if err != nil {
		return "", err
	}

	return "", err
}
