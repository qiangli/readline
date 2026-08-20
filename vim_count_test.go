package readline

import (
	"io"
	"strings"
	"testing"
)

func TestVimCountedBackwardMotion(t *testing.T) {
	rl, err := NewFromConfig(&Config{
		Stdin:          strings.NewReader("abcd\x1b4hrZ\n"),
		Stdout:         io.Discard,
		Stderr:         io.Discard,
		VimMode:        true,
		FuncIsTerminal: func() bool { return false },
	})
	if err != nil {
		t.Fatal(err)
	}
	defer rl.Close()

	got, err := rl.ReadLine()
	if err != nil {
		t.Fatal(err)
	}
	if want := "Zbcd"; got != want {
		t.Fatalf("ReadLine() = %q, want %q", got, want)
	}
}
