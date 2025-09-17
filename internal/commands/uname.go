package commands

import (
	"fmt"
	"github.com/temkinsx/config-management-REPL/internal/commands/model"
	"os"
	"runtime"
)

type Uname struct{}

func (u *Uname) Name() string {
	return "uname"
}

func (u *Uname) Run(args []string, env *model.Env) (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("error: can't get hostname: %s", err)
	}
	return fmt.Sprintf("%v %v %v", runtime.GOOS, hostname, runtime.GOARCH), nil
}
