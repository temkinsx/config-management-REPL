package commands

import "github.com/temkinsx/config-management-REPL/internal/commands/model"

type Cat struct{}

func (c *Cat) Name() string {
	return "cat"
}

func (c *Cat) Run(args []string, env *model.Env) (string, error) {
	path := args[0]
	if path == "" {
		return "", nil
	}

	node, err := env.FS.ResolvePath(path)
	if err != nil {
		return "", err
	}

	return node.ContentText, nil
}
