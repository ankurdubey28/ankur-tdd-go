package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

type config struct {
	files         []string
	pattern       string
	caseSensitive bool
}

var ErrIsDirectory = errors.New("Directory given")
var ErrFileNotFound = errors.New("file does not exist")
var ErrPermDenied = errors.New("permission to read file denied")

func main() {
	c, err := parseArgs(os.Stdout, os.Args[1:])
	if err != nil {
		fmt.Println("error")
		os.Exit(1)
	}
	err = runCmd(&c)
	if err != nil {
		fmt.Println("error")
		os.Exit(1)
	}
}

func parseArgs(w io.Writer, args []string) (config, error) {
	c := config{}
	fs := flag.NewFlagSet("grep", flag.ContinueOnError)

	fs.Usage = func() {
		fmt.Fprintf(w, "usage: %s <pattern> [options] [filename]\n", fs.Name())
		fmt.Fprintf(w, "searches for pattern in given file and print entire line")
	}

	fs.BoolVar(&c.caseSensitive, "i", false, "ignore case sensitivity")

	err := fs.Parse(args)

	if err != nil {
		return c, err
	}

	input := fs.Args()
	c.pattern = input[0]
	c.files = input[1:]
	return c, nil
}

func validateArgs(file string) error {
	info, err := os.Stat(file)
	if err != nil {
		return ErrFileNotFound
	}
	if info.IsDir() {
		return ErrIsDirectory
	}
	mode := info.Mode()
	if mode.Perm()&(1<<8) == 0 {
		return ErrPermDenied
	}
	return nil
}

func runCmd(c *config) error {
	re, err := regexp.Compile(c.pattern)
	if err != nil {
		return err
	}

	multiFile := len(c.files) > 1

	for _, fname := range c.files {
		if err := validateArgs(fname); err != nil {
			fmt.Fprintln(os.Stderr, "grep:", fname+":", err)
			continue
		}

		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "grep:", fname+":", err)
			continue
		}
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			line := scanner.Text()
			if re.MatchString(line) {
				if multiFile {
					fmt.Printf("%s:%s\n", fname, line)
				} else {
					fmt.Println(line)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "grep:", fname+":", err)
		}

		if err := f.Close(); err != nil {
			fmt.Fprintln(os.Stderr, "grep:", fname+":", err)
		}
	}

	return nil
}
