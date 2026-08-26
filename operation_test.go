package readline

import (
	"bytes"
	"strings"
	"testing"
)

// startupProbeBytes returns the bytes operation.Runes writes to the terminal
// before printing the prompt, when the cursor-position (DSR/CPR) probe runs:
// a " \b" pad (see getAndSetOffset) followed by the "\x1b[6n" DSR query.
func startupProbeBytes(t *testing.T, cfg *Config) []byte {
	t.Helper()
	if err := cfg.init(); err != nil {
		t.Fatal(err)
	}
	rl, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer rl.Close()

	got, err := rl.ReadLine()
	if err != nil {
		t.Fatalf("ReadLine() error = %v", err)
	}
	if got != "" {
		t.Fatalf("ReadLine() = %q, want empty", got)
	}
	return cfg.Stdout.(*bytes.Buffer).Bytes()
}

func TestCursorPositionQueryDisabled(t *testing.T) {
	cfg := &Config{
		Stdin:                      strings.NewReader("\n"),
		Stdout:                     &bytes.Buffer{},
		Stderr:                     &bytes.Buffer{},
		FuncGetSize:                func() (int, int) { return 80, 24 },
		FuncIsTerminal:             func() bool { return true },
		DisableCursorPositionQuery: true,
	}
	out := startupProbeBytes(t, cfg)

	if bytes.Contains(out, []byte("\x1b[6n")) {
		t.Fatalf("DisableCursorPositionQuery: output contains DSR query: %q", out)
	}
	// getAndSetOffset always writes its " \b" line-edge pad as the very first
	// bytes of Runes(), immediately before the DSR query; a bare " \b" can
	// also occur later as part of unrelated prompt refresh, so a prefix
	// check is what specifically distinguishes the probe's own padding.
	if bytes.HasPrefix(out, []byte(" \b")) {
		t.Fatalf("DisableCursorPositionQuery: output starts with probe padding bytes: %q", out)
	}
}

func TestCursorPositionQueryEnabledByDefault(t *testing.T) {
	cfg := &Config{
		// A CPR response for the DSR query, followed by Enter to finish the read.
		Stdin:          strings.NewReader("\x1b[1;1R\n"),
		Stdout:         &bytes.Buffer{},
		Stderr:         &bytes.Buffer{},
		FuncGetSize:    func() (int, int) { return 80, 24 },
		FuncIsTerminal: func() bool { return true },
	}
	out := startupProbeBytes(t, cfg)

	if !bytes.Contains(out, []byte("\x1b[6n")) {
		t.Fatalf("default config: output missing DSR query: %q", out)
	}
	if !bytes.HasPrefix(out, []byte(" \b")) {
		t.Fatalf("default config: output missing probe padding bytes: %q", out)
	}
}

func TestPasswordConfigPreservesDisabledCursorPositionQuery(t *testing.T) {
	cfg := &Config{
		Stdin:                      strings.NewReader("\n"),
		Stdout:                     &bytes.Buffer{},
		Stderr:                     &bytes.Buffer{},
		FuncGetSize:                func() (int, int) { return 80, 24 },
		FuncIsTerminal:             func() bool { return true },
		DisableCursorPositionQuery: true,
	}
	if err := cfg.init(); err != nil {
		t.Fatal(err)
	}
	rl, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer rl.Close()
	if got := rl.GeneratePasswordConfig().DisableCursorPositionQuery; !got {
		t.Fatal("password config re-enabled the cursor-position query")
	}
}
