package readline

import (
	"os"
	"os/exec"
	"strings"
)

const (
	vim_NORMAL = iota
	vim_INSERT
	vim_VISUAL
)

type opVim struct {
	op      *operation
	vimMode int
}

func newVimMode(op *operation) *opVim {
	ov := &opVim{
		op:      op,
		vimMode: vim_INSERT,
	}
	return ov
}

func (o *opVim) IsEnableVimMode() bool {
	return o.op.GetConfig().VimMode
}

func (o *opVim) handleVimNormalMovement(r rune, readNext func() rune) (t rune, handled bool) {
	rb := o.op.buf
	handled = true
	switch r {
	case 'h':
		t = CharBackward
	case 'j':
		t = CharNext
	case 'k':
		t = CharPrev
	case 'l':
		t = CharForward
	case '0', '^':
		rb.MoveToLineStart()
	case '$':
		rb.MoveToLineEnd()
	case 'x':
		rb.Delete()
		if rb.IsCursorInEnd() {
			rb.MoveBackward()
		}
	case 'r':
		rb.Replace(readNext())
	case 'd':
		next := readNext()
		switch next {
		case 'd':
			rb.Erase()
		case 'w':
			rb.DeleteWord()
		case 'h':
			rb.Backspace()
		case 'l':
			rb.Delete()
		}
	case 'p':
		rb.Yank()
	case 'b', 'B':
		rb.MoveToPrevWord()
	case 'w', 'W':
		rb.MoveToNextWord()
	case 'e', 'E':
		rb.MoveToEndWord()
	case 'f', 'F', 't', 'T':
		next := readNext()
		prevChar := r == 't' || r == 'T'
		reverse := r == 'F' || r == 'T'
		switch next {
		case CharEsc:
		default:
			rb.MoveTo(next, prevChar, reverse)
		}
	default:
		return r, false
	}
	return t, true
}

func (o *opVim) handleVimNormalEnterInsert(r rune, readNext func() rune) (t rune, handled bool) {
	rb := o.op.buf
	handled = true
	switch r {
	case 'i':
	case 'I':
		rb.MoveToLineStart()
	case 'a':
		rb.MoveForward()
	case 'A':
		rb.MoveToLineEnd()
	case 's':
		rb.Delete()
	case 'S':
		rb.Erase()
	case 'c':
		next := readNext()
		switch next {
		case 'c':
			rb.Erase()
		case 'w':
			rb.DeleteWord()
		case 'h':
			rb.Backspace()
		case 'l':
			rb.Delete()
		}
	default:
		return r, false
	}

	o.EnterVimInsertMode()
	return
}

func (o *opVim) HandleVimNormal(r rune, readNext func() rune) (t rune) {
	switch r {
	case CharEnter, CharInterrupt:
		o.vimMode = vim_INSERT // ???
		return r
	case 'v':
		// vi `v` command: edit the current line in an external editor, then
		// (bash semantics) accept/execute the edited line. POSIX behavior #28:
		// in POSIX mode the editor is `vi` directly; otherwise $VISUAL, then
		// $EDITOR, then `vi`.
		return o.editAndExecute()
	}

	if r, handled := o.handleVimNormalMovement(r, readNext); handled {
		return r
	}

	if r, handled := o.handleVimNormalEnterInsert(r, readNext); handled {
		return r
	}

	// invalid operation
	o.op.t.Bell()
	return 0
}

// vimEditorCommand returns the editor `v` should invoke, per POSIX bash:
//   - POSIX mode (POSIXLY_CORRECT set): `vi` directly, ignoring $VISUAL/$EDITOR.
//   - otherwise: $VISUAL, then $EDITOR, then `vi`.
func vimEditorCommand() string {
	if os.Getenv("POSIXLY_CORRECT") != "" {
		return "vi"
	}
	if e := os.Getenv("VISUAL"); e != "" {
		return e
	}
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	return "vi"
}

// editAndExecute implements the vi-mode `v` command: it writes the current
// line to a temp file, opens it in the editor (see vimEditorCommand), reads the
// result back as the new line, and returns CharEnter so the shell accepts and
// executes it — matching bash's vi `v` (edit-and-execute) behavior. On any
// failure it rings the bell and stays in normal mode (returns 0).
func (o *opVim) editAndExecute() rune {
	tmp, err := os.CreateTemp("", "bashy-fc-*.sh")
	if err != nil {
		o.op.t.Bell()
		return 0
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err := tmp.WriteString(string(o.op.buf.Runes())); err != nil {
		_ = tmp.Close()
		o.op.t.Bell()
		return 0
	}
	_ = tmp.Close()

	editor := vimEditorCommand()
	// Run the editor with the shell's terminal in cooked mode. The editor is a
	// full-screen program; it needs the real TTY, so wire stdio to the shell's
	// configured streams (default os.Stdin/out/err) and suspend readline's raw
	// mode for the duration.
	cfg := o.op.GetConfig()
	args := append(editorArgs(editor), name)
	cmd := exec.Command(shellWord(editor), args...) //nolint:gosec // user's configured editor
	cmd.Stdin = cfg.Stdin
	cmd.Stdout = cfg.Stdout
	cmd.Stderr = cfg.Stderr
	_ = o.op.t.ExitRawMode()
	runErr := cmd.Run()
	_ = o.op.t.EnterRawMode()
	if runErr != nil {
		o.op.t.Bell()
		return 0
	}

	edited, err := os.ReadFile(name)
	if err != nil {
		o.op.t.Bell()
		return 0
	}
	// Bash feeds the edited buffer back a line at a time and executes it; we
	// take the file content verbatim minus a single trailing newline.
	text := strings.TrimSuffix(string(edited), "\n")
	o.op.SetBuffer(text)
	o.vimMode = vim_INSERT
	return CharEnter
}

// shellWord / editorArgs split an EDITOR value that may carry arguments
// (e.g. "code -w", "emacsclient -nw") into the program and its leading args.
func shellWord(editor string) string {
	fields := strings.Fields(editor)
	if len(fields) == 0 {
		return editor
	}
	return fields[0]
}

func editorArgs(editor string) []string {
	fields := strings.Fields(editor)
	if len(fields) <= 1 {
		return nil
	}
	return fields[1:]
}

func (o *opVim) EnterVimInsertMode() {
	o.vimMode = vim_INSERT
}

func (o *opVim) ExitVimInsertMode() {
	o.vimMode = vim_NORMAL
}

func (o *opVim) HandleVim(r rune, readNext func() rune) rune {
	if o.vimMode == vim_NORMAL {
		return o.HandleVimNormal(r, readNext)
	}
	if r == CharEsc {
		o.ExitVimInsertMode()
		return 0
	}

	switch o.vimMode {
	case vim_INSERT:
		return r
	case vim_VISUAL:
	}
	return r
}
