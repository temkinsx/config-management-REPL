package commands

import (
	"fmt"
	"github.com/temkinsx/config-management-REPL/internal/commands/model"
	"strings"
)

type Ls struct{}

func (l *Ls) Name() string {
	return "ls"
}

func (l *Ls) Run(args []string, env *model.Env) (string, error) {
	if len(args) == 1 {
		target, err := env.FS.ResolvePath(args[0])
		if err != nil {
			return "", err
		}

		if target.Type != "dir" {
			return "", fmt.Errorf("%s is not a directory", target.Name)
		}

		var out []string
		for _, ch := range target.Children {
			out = append(out, ch.Name)
		}

		return strings.Join(out, "\n"), nil
	}

	items, err := env.FS.List()
	if err != nil {
		return "", err
	}

	var out []string
	for _, ch := range items {
		out = append(out, ch.Name)
	}

	return strings.Join(out, "\n"), nil
}
