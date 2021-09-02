
setup: install-pre-commit-hooks
	go get -u -a golang.org/x/tools/cmd/stringer
	go mod tidy

install-pre-commit-hooks:
	pre-commit install --install-hooks

run-pre-commit:
	pre-commit run -a
