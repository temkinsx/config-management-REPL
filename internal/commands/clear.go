package commands

import (
	"fmt"
	"github.com/temkinsx/config-management-REPL/internal/commands/model"
)

type Clear struct{}

func (c *Clear) Name() string {
	return "clear"
}

func (c *Clear) Run(args []string, env *model.Env) (string, error) {
	// ANSI-коды очистки экрана + перемещение курсора в (0,0)
	fmt.Println("\033[2J\033[H")
	return "", nil
}
