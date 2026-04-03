package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type config struct {
	fileName string
	pattern  string
}

var ErrIsDirectory = errors.New("Directory given")
var ErrFileNotFound = errors.New("file does not exist")

func main() {
	runCmd(&config{
		fileName: "C:\\Users\\ASUS\\GolangProjects\\ankur-tdd-go\\cli\\grep\\text.txt",
		pattern:  "hello",
	})
}

func parseArgs(w io.Writer, args []string) (config, error) {
	c := config{}
	fs := flag.NewFlagSet("grep", flag.ContinueOnError)

	fs.Usage = func() {
		fmt.Fprintf(w, "usage: %s <pattern> [options] [filename]\n", fs.Name())
		fmt.Fprintf(w, "searches for pattern in given file and print entire line")
	}

	err := fs.Parse(args)

	if err != nil {
		return c, err
	}

	posArgs := fs.Args()
	c.pattern = posArgs[0]
	c.fileName = posArgs[1]
	return c, nil
}

func validateArgs(c *config) error {
	stats, err := os.Stat(c.fileName)
	if err != nil {
		return ErrFileNotFound
	}
	if stats.IsDir() {
		return ErrIsDirectory
	}
	return nil
}

func runCmd(c *config) error {
	file, err := os.Open(c.fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		content := scanner.Text()
		if strings.Contains(content, c.pattern) {
			fmt.Println(content)
		}
	}
	return nil
}
