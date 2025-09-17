package model

import (
	"github.com/temkinsx/config-management-REPL/internal/vfs"
)

type Env struct {
	FS *vfs.VFS
}

type Command interface {
	Name() string
	Run(args []string, env *Env) (string, error)
}
