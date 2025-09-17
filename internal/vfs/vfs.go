package vfs

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type VFS struct {
	Root *Node
	Cwd  *Node
}

type Node struct {
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	ContentText   string  `json:"contentText,omitempty"`
	ContentBase64 string  `json:"contentBase64,omitempty"`
	Children      []*Node `json:"children,omitempty"`

	parent *Node `json:"-"`
}

func LoadFS(path string) (*VFS, error) {
	var root *Node
	if path != "" {
		f, err := os.Open(path)
		if err != nil {
			return &VFS{}, err
		}
		defer f.Close()

		dec := json.NewDecoder(f)
		var tmp Node
		err = dec.Decode(&tmp)
		if err != nil {
			return &VFS{}, err
		}

		if tmp.Name == "" {
			tmp.Name = "/"
		}

		err = setParents(&tmp, nil)
		if err != nil {
			return &VFS{}, err
		}
		root = &tmp
	} else {
		fmt.Println("Path for JSON not set: using default vfs...")
		root = &Node{
			Name:     "/",
			Type:     "dir",
			Children: []*Node{},
		}
	}

	return &VFS{
		Root: root,
		Cwd:  root,
	}, nil
}

func (v *VFS) ResolvePath(path string) (*Node, error) {
	path = strings.TrimSpace(path)

	if path == "" {
		return v.Cwd, nil
	}

	var start *Node
	if strings.HasPrefix(path, "/") || strings.HasPrefix(path, "./") {
		start = v.Root
	} else {
		start = v.Cwd
	}

	parts := strings.FieldsFunc(path, func(r rune) bool {
		return r == '/'
	})

	cur := start
	for _, part := range parts {
		switch part {
		case "", ".":
			continue
		case "..":
			if cur.parent != nil {
				cur = cur.parent
			}
		default:
			found := false
			for _, ch := range cur.Children {
				if ch.Name == part {
					cur = ch
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("no such file or directory: %s", part)
			}
		}
	}

	return cur, nil
}

func (v *VFS) Touch(path string, name string) error {
	targetDir, err := v.ResolvePath(path)
	if err != nil {
		return err
	}

	if targetDir.Type != "dir" {
		return fmt.Errorf("%s is not a directory", targetDir.Name)
	}

	for _, ch := range targetDir.Children {
		if ch.Name == name {
			return nil
		}
	}

	file := &Node{
		Name:   name,
		Type:   "file",
		parent: targetDir,
	}

	if targetDir.Children == nil {
		targetDir.Children = []*Node{}
	}

	targetDir.Children = append(targetDir.Children, file)

	return nil
}

func (v *VFS) List() ([]*Node, error) {
	if v.Cwd.Type != "dir" {
		return []*Node{}, fmt.Errorf("%s is not a directory", v.Cwd.Name)
	}
	if v.Cwd.Children == nil {
		return []*Node{}, nil
	}
	return v.Cwd.Children, nil
}

func (v *VFS) Cd(path string) error {
	cur, err := v.ResolvePath(path)
	if err != nil {
		return err
	}

	if cur.Type != "dir" {
		return fmt.Errorf("%s is not a directory")
	}

	v.Cwd = cur
	return nil
}

func setParents(n *Node, parent *Node) error {
	n.parent = parent

	switch n.Type {
	case "file":
		if len(n.Children) > 0 {
			return fmt.Errorf("file %s can't have children", n.Name)
		}
	case "dir":
		if len(n.Children) == 0 {
			n.Children = []*Node{}
		}

		seen := make(map[string]struct{})

		//проверка на уникальность имен содержимого
		for _, ch := range n.Children {
			if _, ok := seen[ch.Name]; ok {
				return fmt.Errorf("duplicate child name %q in directory %q", ch.Name, n.Name)
			}
			seen[ch.Name] = struct{}{}
		}
	default:
		return fmt.Errorf("invalid type: %s", n.Type)
	}

	for _, ch := range n.Children {
		err := setParents(ch, n)
		if err != nil {
			return err
		}
	}

	return nil
}

func AbsPath(n *Node) string {
	if n == nil {
		return "/"
	}

	var parts []string
	cur := n
	for cur != nil && cur.parent != nil {
		parts = append(parts, cur.Name)
		cur = cur.parent
	}

	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}

	if len(parts) == 0 {
		return "/"
	}

	return "/" + strings.Join(parts, "/")
}

func PrintVFS(n *Node, indent string) {
	fmt.Printf("%s- %s (%s)\n", indent, n.Name, n.Type)
	for _, child := range n.Children {
		PrintVFS(child, indent+"  ")
	}
}
