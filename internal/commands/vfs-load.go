package commands

import (
	"fmt"
	"github.com/temkinsx/config-management-REPL/internal/commands/model"
	"github.com/temkinsx/config-management-REPL/internal/vfs"
)

type VFSLoad struct{}

func (V *VFSLoad) Name() string {
	return "vfs-load"
}

func (V *VFSLoad) Run(args []string, env *model.Env) (string, error) {
	if len(args) != 0 {
		return "", fmt.Errorf("vfs-load: invalid arguments (vfs-load [path])")
	}

	path := args[0]

	newFS, err := vfs.LoadFS(path)
	if err != nil {
		return "", err
	}

	env.FS = newFS
	return fmt.Sprintf("Loaded new file system: %s", path), nil
}
