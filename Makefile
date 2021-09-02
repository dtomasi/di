
setup: install-pre-commit-hooks
	go mod tidy

install-pre-commit-hooks:
	pre-commit install --install-hooks

run-pre-commit:
	pre-commit run -a
