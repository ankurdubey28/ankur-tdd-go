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
	nBool, cBool, bBool, wBool bool

	fileName string
}

func main() {
	c, err := parseArgs(os.Stdout, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
	// get filename from args
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, err)
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

}

func parseArgs(w io.Writer, args []string) (config, error) {
	c := config{}
	fs := flag.NewFlagSet("wc", flag.ContinueOnError)
	fs.SetOutput(w)

	fs.Usage = func() {
		fmt.Fprintf(w, "Usage: %s [options] [file]\n", fs.Name())
		fmt.Fprintln(w, "Counts lines, words, bytes, and characters (like Unix wc).")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: -l , -w , -m , -c")
		fs.PrintDefaults()
	}

	// flags defined
	flag.BoolVar(&c.nBool, "l", false, "flag for newline count")
	flag.BoolVar(&c.cBool, "m", false, "flag for character count")
	flag.BoolVar(&c.wBool, "w", false, "flag for word count")
	flag.BoolVar(&c.bBool, "c", false, "flag for byte count")

	err := fs.Parse(args)
	if err != nil {
		return c, err
	}

	if fs.NArg() > 1 {
		return c, errors.New("positional arguments specified")
	}
	if fs.NArg() == 1 {
		c.fileName = fs.Arg(0)
	}

	return c, nil
}

func validateArgs(c config) error {
	if c.fileName == "" {
		return errors.New("empty file")
	}
	stats, err := os.Stat(c.fileName)
	if err != nil {
		return err
	}
	if stats.IsDir() {
		return errors.New("directory entered")
	}
	return nil
}

func runCmd(c *config) error {
	lines, words, bytes, chars, err := wcCount(c.fileName)
	if err != nil {
		return err
	}

	// default: -l -w -c
	if !c.nBool && !c.wBool && !c.bBool && !c.cBool {
		c.nBool = true
		c.wBool = true
		c.bBool = true
	}
	// collect outputs in fixed wc order
	var outputs []int

	if c.nBool {
		outputs = append(outputs, lines)
	}
	if c.wBool {
		outputs = append(outputs, words)
	}
	if c.bBool {
		outputs = append(outputs, bytes)
	}
	if c.cBool {
		outputs = append(outputs, chars)
	}
	for _, v := range outputs {
		fmt.Print(v, " ")
	}
	fmt.Println(c.fileName)
	return nil
}

func wcCount(fileName string) (int, int, int, int, error) {
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
	return lines + 1, words, bytes, chars, nil
}
