package commands

import (
	"github.com/temkinsx/config-management-REPL/internal/commands/model"
	"os"
)

type Exit struct{}

func (e *Exit) Name() string {
	return "exit"
}

func (e *Exit) Run(args []string, env *model.Env) (string, error) {
	os.Exit(1)
	return "", nil
}
