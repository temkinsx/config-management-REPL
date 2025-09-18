package commands

import (
	"fmt"
	"strings"

	"github.com/temkinsx/config-management-REPL/internal/commands/model"
)

type Uniq struct {
}

func (u Uniq) Name() string {
	return "uniq"
}

func (u Uniq) Run(args []string, env *model.Env) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("error: invalid arguments\nusage: uniq [pathToFile]")
	}

	targetFile, err := env.FS.ResolvePath(args[0])
	if err != nil {
		return "", err
	}

	if targetFile.Type != "file" {
		return "", fmt.Errorf("error: %s is not a file", targetFile)
	}

	lines := strings.Split(targetFile.ContentText, "\n")

	if len(lines) < 1 {
		return "", nil
	}

	var result []string
	var prev string
	for i, s := range lines {
		if i == 0 || s != prev {
			result = append(result, s)
		}
		prev = s
	}

	targetFile.ContentText = strings.Join(result, "\n")
	return targetFile.ContentText, nil
}
