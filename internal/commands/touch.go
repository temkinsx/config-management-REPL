package commands

import (
	"fmt"
	"github.com/temkinsx/config-management-REPL/internal/commands/model"
	"strings"
)

type Touch struct{}

func (t *Touch) Name() string {
	return "touch"
}

func (t *Touch) Run(args []string, env *model.Env) (string, error) {
	switch len(args) {
	case 0:
		return "", nil
	case 1:
		targetDir, name := getNameFromPath(args[0])
		err := env.FS.Touch(targetDir, name)
		if err != nil {
			return "", err
		}
		return "", err
	default:
		return "", fmt.Errorf("touch: invalid arguments")
	}
}

func getNameFromPath(path string) (string, string) {
	parts := strings.Split(path, "/")
	name := parts[len(parts)-1]

	targetDir, ok := strings.CutSuffix(path, name)
	if !ok {
		return "", ""
	}

	return targetDir, name
}
