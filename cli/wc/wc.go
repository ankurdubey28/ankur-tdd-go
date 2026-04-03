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

var IsDirectory = errors.New("Is Directory")

type config struct {
	nBool, cBool, bBool, wBool bool

	fileNames []string
}

func main() {
	c, err := parseArgs(os.Stdout, os.Args[1:])
	if err != nil {
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
	fs.BoolVar(&c.nBool, "l", false, "flag for newline count")
	fs.BoolVar(&c.cBool, "m", false, "flag for character count")
	fs.BoolVar(&c.wBool, "w", false, "flag for word count")
	fs.BoolVar(&c.bBool, "c", false, "flag for byte count")

	err := fs.Parse(args)
	if err != nil {
		return c, err
	}

	c.fileNames = fs.Args()

	return c, nil
}

func validateArgs(fileName string) error {
	stats, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	if stats.IsDir() {
		return IsDirectory
	}
	return nil
}

func runCmd(c *config) error {

	if len(c.fileNames) == 0 {
		lines, words, bytes, chars, err := wcFromReader(os.Stdin)
		if err != nil {
			return err
		}

		if !c.nBool && !c.wBool && !c.bBool && !c.cBool {
			c.nBool = true
			c.wBool = true
			c.bBool = true
		}

		if c.nBool {
			fmt.Print(lines, " ")
		}
		if c.wBool {
			fmt.Print(words, " ")
		}
		if c.bBool {
			fmt.Print(bytes, " ")
		}
		if c.cBool {
			fmt.Print(chars, " ")
		}
		fmt.Println()

		return nil
	}

	var totalCount = []int{0, 0, 0, 0}

	if !c.nBool && !c.wBool && !c.bBool && !c.cBool {
		c.nBool = true
		c.wBool = true
		c.bBool = true
	}

	for _, fileName := range c.fileNames {
		lines, words, bytes, chars := 0, 0, 0, 0
		hasError := false

		if err := validateArgs(fileName); err != nil {
			hasError = true

			if errors.Is(err, IsDirectory) {
				fmt.Fprintf(os.Stderr, "wc: %s: Is a directory\n", fileName)

			} else {
				fmt.Fprintf(os.Stderr, "wc: %s: %v\n", fileName, err)
				continue
			}
		}

		if !hasError {
			file, err := os.Open(fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "wc: %s: %v\n", fileName, err)
				continue
			}

			lines, words, bytes, chars, err = wcFromReader(file)
			file.Close()

			if err != nil {
				fmt.Fprintf(os.Stderr, "wc: %s: %v\n", fileName, err)
				continue
			}
		}

		// print output
		if c.nBool {
			totalCount[0] += lines
			fmt.Print(lines, " ")
		}
		if c.wBool {
			totalCount[1] += words
			fmt.Print(words, " ")
		}
		if c.bBool {
			totalCount[2] += bytes
			fmt.Print(bytes, " ")
		}
		if c.cBool {
			totalCount[3] += chars
			fmt.Print(chars, " ")
		}
		fmt.Println(fileName)
	}

	if len(c.fileNames) > 1 {
		if c.nBool {
			fmt.Print(totalCount[0], " ")
		}
		if c.wBool {
			fmt.Print(totalCount[1], " ")
		}
		if c.bBool {
			fmt.Print(totalCount[2], " ")
		}
		if c.cBool {
			fmt.Print(totalCount[3], " ")
		}
		fmt.Println("total")
	}
	return nil
}

func wcFromReader(r io.Reader) (int, int, int, int, error) {
	reader := bufio.NewReader(r)

	var lines, words, bytes, chars int
	inWord := false

	for {
		r, size, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, 0, 0, 0, err
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
