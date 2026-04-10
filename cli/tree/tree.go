package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	args := []string{"."}

	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	for _, arg := range args {
		fmt.Println(arg)
		err := tree(arg, "")
		if err != nil {
			log.Printf("tree %s: %v\n", arg, err)
		}
	}
}

func tree(root, indent string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return fmt.Errorf("could not read %s: %v", root, err)
	}
	// filter hidden files
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
			return err
		}

		if info.IsDir() {
			if err := tree(fullPath, indent+nextIndent); err != nil {
				return err
			}
		}
	}

	return nil
}
