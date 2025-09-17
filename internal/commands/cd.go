package commands

import (
	"fmt"
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

	targetDir, err := env.FS.ResolvePath(args[0])
	if err != nil {
		return "", err
	}
	if targetDir.Type != "dir" {
		return "", fmt.Errorf("%s is not a directory")
	}

	env.FS.Cwd = targetDir
	return "", err
}
