package readline

import "testing"

// TestVimEditorCommand verifies the POSIX behavior #28 editor selection for the
// vi-mode `v` (edit-and-execute) command: POSIX mode invokes `vi` directly,
// ignoring $VISUAL/$EDITOR; otherwise $VISUAL, then $EDITOR, then `vi`.
// This is the conformance-relevant kernel of #28 (the interactive trigger is
// terminal-dependent and out of scope for automated conformance).
func TestVimEditorCommand(t *testing.T) {
	cases := []struct {
		name                       string
		posixlyCorrect, visual, ed string
		want                       string
	}{
		{"posix ignores VISUAL/EDITOR", "y", "code", "emacs", "vi"},
		{"posix with nothing set", "1", "", "", "vi"},
		{"non-posix prefers VISUAL", "", "nvim", "emacs", "nvim"},
		{"non-posix falls to EDITOR", "", "", "emacs", "emacs"},
		{"non-posix default vi", "", "", "", "vi"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("POSIXLY_CORRECT", tc.posixlyCorrect)
			t.Setenv("VISUAL", tc.visual)
			t.Setenv("EDITOR", tc.ed)
			if got := vimEditorCommand(); got != tc.want {
				t.Fatalf("vimEditorCommand() = %q, want %q", got, tc.want)
			}
		})
	}
}
