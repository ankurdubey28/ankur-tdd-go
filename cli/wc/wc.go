package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode"
)

type config struct {
	newLineCount int
	charCount    int
	byteCount    int
	wordCount    int
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
	fmt.Println(c.newLineCount, c.wordCount, c.charCount, c.byteCount, c.fileName)
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
	lines, words, bytes, chars, err := wc(c.fileName)

	if err != nil {
		return errors.New("error evaluating file")
	}

	c.newLineCount = lines
	c.charCount = chars
	c.byteCount = bytes
	c.wordCount = words
	return nil
}

func wc(fileName string) (int, int, int, int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var lines, words, bytes, chars int
	inWord := false

	for {
		r, size, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		// bytes
		bytes += size
		// characters (runes)
		chars++
		// lines
		if r == '\n' {
			lines++
		}
		// words
		if unicode.IsSpace(r) {
			inWord = false
		} else if !inWord {
			words++
			inWord = true
		}
	}
	return lines, words, bytes, chars, nil
}
