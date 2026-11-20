###########################################################################
##  GENERATED FILE — DO NOT EDIT.                                        ##
##  Produced from dag.md by `bashy dag --makefile`. Edit dag.md instead, ##
##  then regenerate; manual changes are overwritten and fail --check.    ##
###########################################################################
# Tasks live in dag.md (a bashy dag file); this Makefile only forwards.
#   bashy dag <target>           run a target (its dependencies first)
#   bashy dag --list             list targets
#   bashy dag --dry-run <target> print the plan, run nothing
#   bashy dag --skill DIR        export these tasks as an Agent Skill
# Regenerate with: bashy dag --makefile  (verify: bashy dag --makefile --check)

# List the dag targets, then how to run them natively.
help: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag --list; echo "Run: bashy dag <target> | skill: bashy dag --skill DIR | more: bashy help"

# Race-enabled tests, vet, and the gofmt check — the old Makefile's test body.
test: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag test

# Build the demo example for linux, windows and darwin (compile check only).
examples: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag examples

# Apply gofmt fixes (.check-gofmt.sh --fix).
gofmt: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag gofmt

# The gofmt gate (reports unformatted files, never rewrites them), plus a check that the generated Makefile still matches this file — regenerate with `bashy dag…
fmtcheck: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag fmtcheck

# Install the pre-push formatting gate.
hooks: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag hooks

# The CI gate: fmtcheck, then test, then the examples build.
ci: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag ci

# The default target: everything CI runs.
all: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag all

.DEFAULT: ; @command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag $@

# FORCE names no file and has no recipe: as the one prerequisite it
# keeps every forward running even when a file or directory on disk
# shares a target's name. Dependency order stays in dag, not here.
FORCE:
