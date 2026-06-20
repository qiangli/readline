# Fork notes — qiangli/readline

Fork of [ergochat/readline](https://github.com/ergochat/readline) (itself a
maintained fork of chzyer/readline), MIT-licensed. Kept close to upstream.

## Divergence from upstream

- **vi `v` edit-and-execute command** (`vim.go`): in vi command mode, `v` now
  opens the current line in an external editor and, on save, accepts/executes
  the edited line — matching GNU bash 5.3's vi-mode `v`. Editor selection per
  POSIX behavior #28: when `POSIXLY_CORRECT` is set the editor is `vi` directly;
  otherwise `$VISUAL`, then `$EDITOR`, then `vi`. Self-contained in `vim.go`
  (no public API change) so downstream consumers that don't need it are
  unaffected.

## Why a fork

The pure-Go Bash 5.3 drop-in (github.com/qiangli/bashy, via github.com/qiangli/sh's
`interactive` package) needs the vi `v` command for POSIX-mode conformance.
Upstream ergochat/readline has vi mode but not the `v` command, and exposes no
hook to add it without modifying the package. Consumed via a go.mod `replace`
of `github.com/ergochat/readline`. The `v`-command patch is a candidate to
upstream as a PR; if accepted, this fork can be retired.
