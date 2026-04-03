package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestWcFromReader(t *testing.T) {
	tests := []struct {
		name  string
		input string
		lines int
		words int
		bytes int
		chars int
	}{
		{
			name:  "empty input",
			input: "",
			lines: 0, words: 0, bytes: 0, chars: 0,
		},
		{
			name:  "only newline",
			input: "\n",
			lines: 1, words: 0, bytes: 1, chars: 1,
		},
		{
			name:  "single word",
			input: "hello",
			lines: 0, words: 1, bytes: 5, chars: 5,
		},
		{
			name:  "two words one line",
			input: "hello world",
			lines: 0, words: 2, bytes: 11, chars: 11,
		},
		{
			name:  "two words with newline",
			input: "hello world\n",
			lines: 1, words: 2, bytes: 12, chars: 12,
		},
		{
			name:  "multiple lines",
			input: "a b\nc d\n",
			lines: 2, words: 4, bytes: 8, chars: 8,
		},
		{
			name:  "multiple spaces",
			input: "a   b",
			lines: 0, words: 2, bytes: 5, chars: 5,
		},
		{
			name:  "leading and trailing spaces",
			input: "   hello world   ",
			lines: 0, words: 2, bytes: 17, chars: 17,
		},
		{
			name:  "tabs and spaces",
			input: "hello\tworld",
			lines: 0, words: 2, bytes: 11, chars: 11,
		},
		{
			name:  "unicode characters",
			input: "hello 🧠",
			lines: 0, words: 2, bytes: 10, chars: 7,
		},
		{
			name:  "unicode with newline",
			input: "🧠\n🧠",
			lines: 1, words: 2, bytes: 9, chars: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)

			l, w, b, c, err := wcFromReader(r)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if l != tt.lines || w != tt.words || b != tt.bytes || c != tt.chars {
				t.Errorf("got (%d,%d,%d,%d), want (%d,%d,%d,%d)",
					l, w, b, c,
					tt.lines, tt.words, tt.bytes, tt.chars,
				)
			}
		})
	}
}

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantCfg config
		wantErr bool
	}{
		{
			name: "no flags no files",
			args: []string{},
			wantCfg: config{
				fileNames: []string{},
			},
		},
		{
			name: "single flag",
			args: []string{"-l"},
			wantCfg: config{
				nBool:     true,
				fileNames: []string{},
			},
		},
		{
			name: "multiple flags",
			args: []string{"-l", "-w", "-c", "-m"},
			wantCfg: config{
				nBool: true, wBool: true, bBool: true, cBool: true,
				fileNames: []string{},
			},
		},
		{
			name: "single file",
			args: []string{"file.txt"},
			wantCfg: config{
				fileNames: []string{"file.txt"},
			},
		},
		{
			name: "flags + file",
			args: []string{"-l", "-w", "file.txt"},
			wantCfg: config{
				nBool: true, wBool: true,
				fileNames: []string{"file.txt"},
			},
		},
		{
			name: "multiple files",
			args: []string{"file1.txt", "file2.txt"},
			wantCfg: config{
				fileNames: []string{"file1.txt", "file2.txt"},
			},
		},
		{
			name:    "invalid flag",
			args:    []string{"-x"},
			wantErr: true,
		},
		{
			name: "flags stop after first non-flag",
			args: []string{"-l", "file1.txt", "-w", "file2.txt"},
			wantCfg: config{
				nBool:     true,
				wBool:     false,
				bBool:     false,
				cBool:     false,
				fileNames: []string{"file1.txt", "-w", "file2.txt"},
			},
		},
		{
			name: "flags before files",
			args: []string{"-l", "-w", "file1.txt", "file2.txt"},
			wantCfg: config{
				nBool: true, wBool: true,
				fileNames: []string{"file1.txt", "file2.txt"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := parseArgs(io.Discard, tt.args)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cfg.nBool != tt.wantCfg.nBool ||
				cfg.wBool != tt.wantCfg.wBool ||
				cfg.bBool != tt.wantCfg.bBool ||
				cfg.cBool != tt.wantCfg.cBool {
				t.Errorf("flags mismatch: got %+v, want %+v", cfg, tt.wantCfg)
			}

			if !reflect.DeepEqual(cfg.fileNames, tt.wantCfg.fileNames) {
				t.Errorf("filenames mismatch: got %v, want %v",
					cfg.fileNames, tt.wantCfg.fileNames)
			}
		})
	}
}
func TestValidateArgs(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() string // returns path
		wantErr  bool
		errIsDir bool
	}{
		{
			name: "valid file",
			setup: func() string {
				f, _ := os.CreateTemp("", "testfile")
				f.Close()
				return f.Name()
			},
			wantErr: false,
		},
		{
			name: "non-existent file",
			setup: func() string {
				return "non_existent_file.txt"
			},
			wantErr: true,
		},
		{
			name: "directory",
			setup: func() string {
				dir, _ := os.MkdirTemp("", "testdir")
				return dir
			},
			wantErr:  true,
			errIsDir: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			defer os.RemoveAll(path)

			err := validateArgs(path)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errIsDir && !errors.Is(err, IsDirectory) {
					t.Fatalf("expected IsDirectory error, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestRunCmdStdin(t *testing.T) {
	oldStdin := os.Stdin
	rPipe, wPipe, _ := os.Pipe()
	os.Stdin = rPipe
	defer func() { os.Stdin = oldStdin }()

	// write input
	wPipe.Write([]byte("hello world\n"))
	wPipe.Close()

	// capture stdout
	oldStdout := os.Stdout
	outR, outW, _ := os.Pipe()
	os.Stdout = outW
	defer func() { os.Stdout = oldStdout }()

	c := config{
		nBool: true,
		wBool: true,
		bBool: true,
	}

	err := runCmd(&c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	outW.Close()

	var buf bytes.Buffer
	buf.ReadFrom(outR)

	got := buf.String()
	want := "1 2 12 \n"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRunCmdSingleFile(t *testing.T) {
	f, _ := os.CreateTemp("", "testfile")
	defer os.Remove(f.Name())

	f.WriteString("hello world\n")
	f.Close()

	// capture stdout
	oldStdout := os.Stdout
	rPipe, wPipe, _ := os.Pipe()
	os.Stdout = wPipe
	defer func() { os.Stdout = oldStdout }()

	c := config{
		nBool:     true,
		wBool:     true,
		fileNames: []string{f.Name()},
	}

	err := runCmd(&c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wPipe.Close()

	var buf bytes.Buffer
	buf.ReadFrom(rPipe)

	got := buf.String()
	want := "1 2 " + f.Name() + "\n"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRunCmdMultipleFiles(t *testing.T) {
	f1, _ := os.CreateTemp("", "f1")
	defer os.Remove(f1.Name())
	f1.WriteString("a b\n")
	f1.Close()

	f2, _ := os.CreateTemp("", "f2")
	defer os.Remove(f2.Name())
	f2.WriteString("c d\n")
	f2.Close()

	oldStdout := os.Stdout
	rPipe, wPipe, _ := os.Pipe()
	os.Stdout = wPipe
	defer func() { os.Stdout = oldStdout }()

	c := config{
		nBool:     true,
		wBool:     true,
		fileNames: []string{f1.Name(), f2.Name()},
	}

	err := runCmd(&c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wPipe.Close()

	var buf bytes.Buffer
	buf.ReadFrom(rPipe)

	got := buf.String()

	want :=
		"1 2 " + f1.Name() + "\n" +
			"1 2 " + f2.Name() + "\n" +
			"2 4 total\n"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
