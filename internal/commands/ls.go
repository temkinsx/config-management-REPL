package commands

import (
	"strings"

	"github.com/temkinsx/config-management-REPL/internal/commands/model"
)

type Ls struct{}

func (l *Ls) Name() string {
	return "ls"
}

func (l *Ls) Run(args []string, env *model.Env) (string, error) {
	if len(args) == 1 {
		orig := env.FS.Cwd

		if err := env.FS.Cd(args[0]); err != nil {
			return "", err
		}
		items, err := env.FS.List()

		env.FS.Cwd = orig
		if err != nil {
			return "", err
		}

		out := make([]string, 0, len(items))
		for _, ch := range items {
			out = append(out, ch.Name)
		}
		return strings.Join(out, "\n"), nil
	}

	items, err := env.FS.List()
	if err != nil {
		return "", err
	}
	out := make([]string, 0, len(items))
	for _, ch := range items {
		out = append(out, ch.Name)
	}
	return strings.Join(out, "\n"), nil
}
