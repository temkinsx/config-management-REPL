package vfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	TypeFile = "file"
	TypeDir  = "dir"
)

var (
	ErrNotFound       = errors.New("no such file or directory")
	ErrNotADir        = errors.New("is not a directory")
	ErrFileNoChildren = errors.New("file can't have children")
	ErrDuplicateName  = errors.New("duplicate child name in directory")
	ErrInvalidType    = errors.New("invalid type")
	ErrInvalidName    = errors.New("invalid name")
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
			return nil, err
		}
		defer f.Close()

		dec := json.NewDecoder(f)
		var tmp Node
		err = dec.Decode(&tmp)
		if err != nil {
			return nil, err
		}

		if tmp.Name == "" {
			tmp.Name = "/"
		}

		err = setParents(&tmp, nil)
		if err != nil {
			return nil, err
		}
		root = &tmp
	} else {
		fmt.Println("Path for JSON is not specified: using default vfs...")
		root = &Node{
			Name:     "/",
			Type:     TypeDir,
			Children: []*Node{},
		}
	}

	return &VFS{
		Root: root,
		Cwd:  root,
	}, nil
}

func (v *VFS) Save(path string) error {
	if path == "" {
		fmt.Println("Path is not specified: snapshot will be stored at /snapshots directory")
		path = snapshotFileName("./snapshots")
	}

	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v.Root); err != nil {
		return err
	}

	fmt.Printf("snapshot saved: %s\n", path)

	return nil
}

func (v *VFS) ResolvePath(path string) (*Node, error) {
	path = strings.TrimSpace(path)

	if path == "" {
		return v.Cwd, nil
	}

	var start *Node
	if strings.HasPrefix(path, "/") {
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
				return nil, fmt.Errorf("%w: %s", ErrNotFound, part)
			}
		}
	}

	return cur, nil
}

func (v *VFS) Touch(path string, name string) error {
	if name == "" || strings.Contains(name, "/") || name == "." || name == ".." {
		return fmt.Errorf("%s: %w", name, ErrInvalidName)
	}

	targetDir, err := v.ResolvePath(path)
	if err != nil {
		return err
	}

	if targetDir.Type != TypeDir {
		return fmt.Errorf("%s %w", targetDir.Name, ErrNotADir)
	}

	for _, ch := range targetDir.Children {
		if ch.Name == name {
			return nil
		}
	}

	file := &Node{
		Name:   name,
		Type:   TypeFile,
		parent: targetDir,
	}

	if targetDir.Children == nil {
		targetDir.Children = []*Node{}
	}

	targetDir.Children = append(targetDir.Children, file)

	return nil
}

func (v *VFS) List() ([]*Node, error) {
	if v.Cwd.Type != TypeDir {
		return nil, fmt.Errorf("%s: %w", v.Cwd.Name, ErrNotADir)
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

	if cur.Type != TypeDir {
		return fmt.Errorf("%s: %w", cur.Name, ErrNotADir)
	}

	v.Cwd = cur
	return nil
}

func setParents(n *Node, parent *Node) error {
	n.parent = parent

	switch n.Type {
	case TypeFile:
		if len(n.Children) > 0 {
			return fmt.Errorf("%s: %w", n.Name, ErrFileNoChildren)
		}
	case TypeDir:
		if n.Children == nil {
			n.Children = []*Node{}
		}

		seen := make(map[string]struct{})

		//проверка на уникальность имен содержимого
		for _, ch := range n.Children {
			if _, ok := seen[ch.Name]; ok {
				return fmt.Errorf("%s: %w (in dir %s)", ch.Name, ErrDuplicateName, n.Name)
			}
			seen[ch.Name] = struct{}{}
		}
	default:
		return fmt.Errorf("%w: %s", ErrInvalidType, n.Type)
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

func snapshotFileName(dir string) string {
	now := time.Now()

	timestamp := now.Format("2006-01-02_15-04-05")

	return filepath.Join(dir, "snapshot_"+timestamp+".json")
}
