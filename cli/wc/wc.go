package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type config struct {
	newLineCount int
	charCount    int
	byteCount    int
	fileName     string
}

func main() {
	c, err := parseArgs(os.Stderr, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
	// get filename from args
	args := os.Args[1:]
	if len(args) == 0 {
		// handle stdin later if needed
		os.Exit(1)
	}
	c.fileName = args[0]
	// validate
	if err := validateArgs(c); err != nil {
		os.Exit(1)
	}
	// run command
	if err := runCmd(&c); err != nil {
		os.Exit(1)
	}
	// print result
	// format similar to wc: lines chars bytes filename
	fmt.Println(c.newLineCount, c.charCount, c.byteCount, c.fileName)
}

func parseArgs(w io.Writer, args []string) (config, error) {
	c := config{}
	fs := flag.NewFlagSet("wc", flag.ContinueOnError)
	fs.SetOutput(w)

	err := fs.Parse(args)
	if err != nil {
		return c, err
	}
	return c, nil
}

func validateArgs(c config) error {
	if c.fileName == "" {
		return errors.New("empty file")
	}
	return nil
}

func runCmd(c *config) error {
	lines, chars, bytes, err := wc(c.fileName)

	if err != nil {
		return errors.New("error evaluating file")
	}

	c.newLineCount = lines
	c.charCount = chars
	c.byteCount = bytes
	return nil
}

func wc(fileName string) (int, int, int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, 0, 0, errors.New("file open error")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	nlines, cCount, bCount := 0, 0, 0
	for scanner.Scan() {
		line := scanner.Text()
		nlines++
		// characters (runes)
		cCount += len([]rune(line))

		// bytes
		bCount += len(line)
	}
	if err = scanner.Err(); err != nil {
		return 0, 0, 0, err
	}
	return nlines, cCount, bCount, nil
}
