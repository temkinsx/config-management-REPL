package prompt

import (
	"fmt"
	"github.com/temkinsx/config-management-REPL/internal/vfs"
	"os"
	"os/user"
)

type Prompt struct {
	User     string
	Hostname string
	FS       *vfs.VFS
}

func New(vfs *vfs.VFS) *Prompt {
	u, _ := user.Current()
	username := ""

	if u != nil {
		username = u.Username
	}
	if username == "" {
		username = os.Getenv("USER")
		if username == "" {
			username = os.Getenv("USERNAME")
		}
	}

	host, _ := os.Hostname()

	return &Prompt{
		User:     username,
		Hostname: host,
		FS:       vfs,
	}
}

func (p *Prompt) Build() string {
	return fmt.Sprintf("%s@%s:%s$ ", p.User, p.Hostname, vfs.AbsPath(p.FS.Cwd))
}
