package commands

import (
	"github.com/temkinsx/config-management-REPL/internal/commands/model"
	"strings"
)

type Echo struct{}

func (c *Echo) Name() string {
	return "echo"
}

func (c *Echo) Run(args []string, env *model.Env) (string, error) {
	line := strings.Join(args, " ")
	return line, nil
}
