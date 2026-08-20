package readline

import (
	"errors"
	"io"
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
