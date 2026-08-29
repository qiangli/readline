package readline

import (
	"errors"
	"io"
	"strings"
	"testing"
	"time"
)

type channelWriter chan []byte

func (w channelWriter) Write(p []byte) (int, error) {
	w <- append([]byte(nil), p...)
	return len(p), nil
}

func TestCursorPositionStopsReadingAtUserInput(t *testing.T) {
	stdin, input := io.Pipe()
	defer stdin.Close()
	defer input.Close()

	stdout := make(channelWriter, 1)
	cfg := &Config{
		Stdin:       stdin,
		Stdout:      &stdout,
		Stderr:      io.Discard,
		FuncGetSize: func() (int, int) { return 80, 24 },
	}
	if err := cfg.init(); err != nil {
		t.Fatal(err)
	}
	term, err := newTerminal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer term.Close()

	done := make(chan error, 1)
	go func() {
		_, err := term.GetCursorPosition(nil)
		done <- err
	}()

	select {
	case got := <-stdout:
		if string(got) != "\x1b[6n" {
			t.Fatalf("cursor-position query = %q", got)
		}
	case <-time.After(time.Second):
		t.Fatal("cursor-position query was not written")
	}
	if _, err := input.Write([]byte("x")); err != nil {
		t.Fatal(err)
	}

	select {
	case err := <-done:
		if !errors.Is(err, invalidCPR) {
			t.Fatalf("GetCursorPosition error = %v, want invalidCPR", err)
		}
	case <-time.After(time.Second):
		t.Fatal("GetCursorPosition kept reading after user input")
	}
	if term.inFlight {
		t.Fatal("cursor-position query left a stdin read in flight")
	}
	r, err := term.GetRune(nil)
	if err != nil {
		t.Fatal(err)
	}
	if r != 'x' {
		t.Fatalf("buffered rune = %q, want x", r)
	}
}

func TestTerminalDoesNotReadAheadPastSubmittedLine(t *testing.T) {
	stdin := strings.NewReader("shell\neditor\n")
	cfg := &Config{Stdin: stdin, Stdout: io.Discard, Stderr: io.Discard}
	if err := cfg.init(); err != nil {
		t.Fatal(err)
	}
	term, err := newTerminal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer term.Close()

	for _, want := range "shell\n" {
		got, err := term.GetRune(nil)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("terminal rune = %q, want %q", got, want)
		}
	}

	// Once readline has accepted the shell command, the foreground command
	// owns stdin. Its editor input must still be present in the shared reader.
	rest, err := io.ReadAll(stdin)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(rest), "editor\n"; got != want {
		t.Fatalf("foreground input = %q, want %q", got, want)
	}
}

func TestExactRuneReaderUTF8AndUnread(t *testing.T) {
	r := &exactRuneReader{r: strings.NewReader("éx")}
	got, size, err := r.ReadRune()
	if err != nil || got != 'é' || size != 2 {
		t.Fatalf("first rune = %q size=%d err=%v", got, size, err)
	}
	if err := r.UnreadRune(); err != nil {
		t.Fatal(err)
	}
	got, size, err = r.ReadRune()
	if err != nil || got != 'é' || size != 2 {
		t.Fatalf("unread rune = %q size=%d err=%v", got, size, err)
	}
	got, size, err = r.ReadRune()
	if err != nil || got != 'x' || size != 1 {
		t.Fatalf("final rune = %q size=%d err=%v", got, size, err)
	}
}
