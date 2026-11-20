---
name: readline
description: Build/test/lint targets for readline as a bashy dag pipeline (agent-first; the Makefile is generated from this file and only forwards here)
type: dag
---

# readline — DAG task file

The Makefile in this directory is GENERATED from this file (`bashy dag --makefile`)
and only forwards to the targets below — do not edit it. This DAG file is the
agent-first task list, runnable with `bashy dag`:

```bash
bashy dag --list            # available targets
bashy dag ci                # fmtcheck + test + examples (what CI runs)
bashy dag --json test       # machine-readable envelope for an agent
```

The targets mirror the old Makefile's commands exactly; `make ci` keeps working
through the generated forwarder.

## Tasks

### test
Race-enabled tests, vet, and the gofmt check — the old Makefile's test body.
Effects: exec,net,write

```bash
set -e
go test -race ./...
go vet ./...
./.check-gofmt.sh
```

### examples
Build the demo example for linux, windows and darwin (compile check only).
Effects: exec,net,write

```bash
set -e
GOOS=linux go build -o /dev/null example/readline-demo/readline-demo.go
GOOS=windows go build -o /dev/null example/readline-demo/readline-demo.go
GOOS=darwin go build -o /dev/null example/readline-demo/readline-demo.go
```

### gofmt
Apply gofmt fixes (.check-gofmt.sh --fix). This APPLIES formatting; the gate
that only reports is `fmtcheck`.
Effects: exec,write

```bash
./.check-gofmt.sh --fix
```

### fmtcheck
The gofmt gate (reports unformatted files, never rewrites them), plus a check
that the generated Makefile still matches this file — regenerate with
`bashy dag --makefile` instead of editing it.
Effects: exec,read

```bash
set -e
./scripts/fmtcheck.sh
bashy dag --makefile --check
```

### hooks
Install the pre-push formatting gate.
Effects: exec,write

```bash
git config core.hooksPath scripts/hooks
echo "hooks installed: core.hooksPath=scripts/hooks"
```

### ci
The CI gate: fmtcheck, then test, then the examples build.
Requires: fmtcheck, test, examples

```bash
:
```

### all
The default target: everything CI runs.
Requires: ci

```bash
:
```
