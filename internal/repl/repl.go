package repl

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/temkinsx/config-management-REPL/internal/commands/model"
	"github.com/temkinsx/config-management-REPL/internal/commands/registry"
	"github.com/temkinsx/config-management-REPL/internal/parser"
	"github.com/temkinsx/config-management-REPL/internal/prompt"
	"github.com/temkinsx/config-management-REPL/internal/vfs"
	"os"
	"strings"
)

var (
	vfsPath    = flag.String("vfs", "", "path to Virtual File System JSON\nIf empty - starting with default VFS (internal/vfs/vfs_default.json)")
	scriptPath = flag.String("script", "", "path to start script")
)

type REPL struct{}

func (r *REPL) Run() {
	flag.Parse()
	cmds := registry.NewRegistry()

	fs, err := vfs.LoadFS(*vfsPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p := prompt.New(fs)
	env := &model.Env{FS: fs}

	var sc *bufio.Scanner
	if *scriptPath != "" {
		f, err := os.Open(*scriptPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		sc = bufio.NewScanner(f)
		runScript(sc, p, env, cmds)
	}

	sc = bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(p.Build())

		if !sc.Scan() {
			break
		}

		line := sc.Text()
		if line == "" {
			continue
		}

		cmdName, args, err := parser.ParseLine(line)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}

		cmd, ok := cmds.Commands[cmdName]
		if !ok {
			fmt.Printf("command not found: %s\n", cmdName)
			continue
		}

		out, err := cmd.Run(args, env)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if out != "" {
			fmt.Println(out)
		}
	}
}

func runScript(sc *bufio.Scanner, p *prompt.Prompt, env *model.Env, cmds *registry.Registry) {
	for {
		if !sc.Scan() {
			break
		}

		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fmt.Print(p.Build())

		fmt.Println(line)

		cmdName, args, err := parser.ParseLine(line)
		if err != nil {
			continue
		}

		cmd, ok := cmds.Commands[cmdName]
		if !ok {
			fmt.Printf("# skipped: command not found: %s\n", cmdName)
			continue
		}

		out, err := cmd.Run(args, env)
		if err != nil {
			fmt.Println("# skipped:", err)
			continue
		}

		if out != "" {
			fmt.Println(out)
		}
	}
}
