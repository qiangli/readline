.PHONY: all test examples ci gofmt fmtcheck hooks

all: ci

ci: fmtcheck test examples

test:
	go test -race ./...
	go vet ./...
	./.check-gofmt.sh

examples:
	GOOS=linux go build -o /dev/null example/readline-demo/readline-demo.go
	GOOS=windows go build -o /dev/null example/readline-demo/readline-demo.go
	GOOS=darwin go build -o /dev/null example/readline-demo/readline-demo.go

gofmt:
	./.check-gofmt.sh --fix

fmtcheck:  ## gofmt gate — reports unformatted files, never rewrites them
	@./scripts/fmtcheck.sh

hooks:  ## install the pre-push formatting gate
	@git config core.hooksPath scripts/hooks
	@echo "hooks installed: core.hooksPath=scripts/hooks"
