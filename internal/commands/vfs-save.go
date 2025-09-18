package commands

import (
	"errors"

	"github.com/temkinsx/config-management-REPL/internal/commands/model"
)

type VfsSave struct{}

func (v *VfsSave) Name() string {
	return "vfs-save"
}

func (v *VfsSave) Run(args []string, env *model.Env) (string, error) {
	var path string
	switch len(args) {
	case 0:
		path = ""
	case 1:
		path = args[0]
	default:
		return "", errors.New("usage: vfs-save [PATH]")
	}

	err := env.FS.Save(path)
	if err != nil {
		return "", err
	}

	return "", err
}
