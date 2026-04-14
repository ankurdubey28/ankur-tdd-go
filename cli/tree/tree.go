package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	inputs    []string
	dirCount  int
	fileCount int
	fullPath  bool
	dirOnly   bool
	levels    int
	perms     bool
}

type node struct {
	name        string
	fullPath    string
	isDir       bool
	indent      string
	isLast      bool
	level       int
	permissions os.FileMode
}

func main() {
	c, err := parseArgs(os.Stdin, os.Args)
	if err != nil {
		return
	}

	if err := runCmd(&c); err != nil {
		os.Exit(1)
	}
}

func parseArgs(w io.Writer, args []string) (config, error) {
	c := config{}

	fs := flag.NewFlagSet("tree", flag.ContinueOnError)
	fs.SetOutput(w)

	fs.Usage = func() {
		fmt.Fprintf(w, "usage: %s [options] [input]\n", fs.Name())
		fmt.Fprintf(w, "prints the directory structure as a tree\n")
	}

	fs.BoolVar(&c.fullPath, "f", false, "print full path")
	fs.BoolVar(&c.dirOnly, "d", false, "print directories only")
	fs.IntVar(&c.levels, "L", -1, "max depth")
	fs.BoolVar(&c.perms, "p", false, "print permissions")

	if err := fs.Parse(args[1:]); err != nil {
		return c, err
	}

	c.inputs = fs.Args()
	if len(c.inputs) == 0 {
		c.inputs = []string{"."}
	}

	return c, nil
}

func runCmd(c *config) error {
	for _, arg := range c.inputs {
		fmt.Println(arg)

		var nodes []node
		dirs, files, err := tree(arg, "", &nodes, 0, c.levels, c.dirOnly)
		if err != nil {
			log.Printf("tree %s: %v\n", arg, err)
			continue
		}

		for _, n := range nodes {
			connector := "├── "
			if n.isLast {
				connector = "└── "
			}

			name := n.name
			if c.fullPath {
				name = n.fullPath
			}

			if c.perms {
				name = fmt.Sprintf("[%s] %s", n.permissions.String(), name)
			}

			fmt.Printf("%s%s%s\n", n.indent, connector, name)
		}

		c.dirCount += dirs
		c.fileCount += files
	}

	fmt.Printf("\n%d directories, %d files\n", c.dirCount, c.fileCount)
	return nil
}

func tree(root, indent string, nodes *[]node, lev int, maxLevel int, dirOnly bool) (int, int, error) {
	fileCount, dirCount := 0, 0

	entries, err := os.ReadDir(root)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: %v", root, err)
	}

	// filter entries first (important for correct tree structure)
	var filtered []os.DirEntry
	for _, e := range entries {
		if len(e.Name()) > 0 && e.Name()[0] == '.' {
			continue
		}
		if dirOnly && !e.IsDir() {
			continue
		}
		filtered = append(filtered, e)
	}

	for i, e := range filtered {
		name := e.Name()
		isLast := i == len(filtered)-1

		fullPath := filepath.Join(root, name)

		info, err := e.Info()
		if err != nil {
			return 0, 0, err
		}

		*nodes = append(*nodes, node{
			name:        name,
			fullPath:    fullPath,
			isDir:       info.IsDir(),
			indent:      indent,
			isLast:      isLast,
			level:       lev,
			permissions: info.Mode(),
		})

		// stop recursion if max depth reached
		if maxLevel != -1 && lev >= maxLevel {
			continue
		}

		var nextIndent string
		if isLast {
			nextIndent = "    "
		} else {
			nextIndent = "│   "
		}

		if info.IsDir() {
			dirCount++

			d, f, err := tree(fullPath, indent+nextIndent, nodes, lev+1, maxLevel, dirOnly)
			if err != nil {
				return 0, 0, err
			}

			dirCount += d
			fileCount += f
		} else {
			fileCount++
		}
	}
	return dirCount, fileCount, nil
}
