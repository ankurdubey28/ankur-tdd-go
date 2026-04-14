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
		dirs, files, err := tree(arg, "")
		if err != nil {
			log.Printf("tree %s: %v\n", arg, err)
			continue
		}

		fmt.Printf("\n%d directories, %d files\n", dirs, files)
	}
}

func tree(root, indent string) (int, int, error) {
	fileCount, dirCount := 0, 0

	entries, err := os.ReadDir(root)
	if err != nil {
		return 0, 0, fmt.Errorf("could not read %s: %v", root, err)
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
