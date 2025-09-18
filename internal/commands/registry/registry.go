package registry

import (
	"github.com/temkinsx/config-management-REPL/internal/commands"
	"github.com/temkinsx/config-management-REPL/internal/commands/model"
)

type Registry struct {
	Commands map[string]model.Command
}

func NewRegistry() *Registry {
	r := &Registry{Commands: make(map[string]model.Command)}
	r.registerAll()
	return r
}

func (r *Registry) registerAll() {
	r.Commands["exit"] = &commands.Exit{}
	r.Commands["echo"] = &commands.Echo{}
	r.Commands["ls"] = &commands.Ls{}
	r.Commands["cd"] = &commands.Cd{}
	r.Commands["touch"] = &commands.Touch{}
	r.Commands["clear"] = &commands.Clear{}
	r.Commands["uname"] = &commands.Uname{}
	r.Commands["uniq"] = &commands.Uniq{}
	r.Commands["cat"] = &commands.Cat{}
	r.Commands["vfs-save"] = &commands.VFSSave{}
	r.Commands["vfs-load"] = &commands.VFSLoad{}

}
