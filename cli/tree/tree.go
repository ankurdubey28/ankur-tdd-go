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
	fullPath  bool
	inputs    []string
	dirCount  int
	fileCount int
}

func main() {

	c, err := parseArgs(os.Stdin, os.Args)
	if err != nil {
		return
	}

	for _, arg := range c.inputs {
		fmt.Println(arg)
		dirs, files, err := tree(arg, "")
		if err != nil {
			log.Printf("tree %s: %v\n", arg, err)
			continue
		}
		c.fileCount += files
		c.dirCount += dirs
	}
	fmt.Printf("\n%d directories, %d files\n", c.dirCount, c.fileCount)
}

func parseArgs(w io.Writer, args []string) (config, error) {
	c := config{}
	fs := flag.NewFlagSet("tree", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(w, "usage: %s [options][input]", fs.Name())
		fmt.Fprintf(w, "prints the entire directory structure in form of tree")
	}
	fs.BoolVar(&c.fullPath, "f", false, "prints full path from root folder")

	if err := fs.Parse(args); err != nil {
		return c, err
	}

	c.inputs = fs.Args()
	return c, nil
}

func tree(root, indent string) (int, int, error) {
	fileCount, dirCount := 0, 0

	entries, err := os.ReadDir(root)
	if err != nil {
		return 0, 0, fmt.Errorf(" %s: [%v]", root, err)
	}

	var names []string
	for _, e := range entries {
		if len(e.Name()) > 0 && e.Name()[0] != '.' {
			names = append(names, e.Name())
		}
	}

	for i, name := range names {
		isLast := i == len(names)-1

		var connector, nextIndent string
		if isLast {
			connector = "└── "
			nextIndent = "    "
		} else {
			connector = "├── "
			nextIndent = "│   "
		}

		fmt.Printf("%s%s%s\n", indent, connector, name)

		fullPath := filepath.Join(root, name)
		info, err := os.Stat(fullPath)
		if err != nil {
			return 0, 0, err
		}

		if info.IsDir() {
			dirCount++

			d, f, err := tree(fullPath, indent+nextIndent)
			if err != nil {
				return 0, 0, err
			}

			dirCount += d
			fileCount += f
		} else {
			fileCount++ // count file
		}
	}
	return dirCount, fileCount, nil
}
