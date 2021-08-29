setup: install-pre-commit-hooks
	go mod tidy

install-pre-commit-hooks:
	pre-commit install --install-hooks

run-pre-commit:
	pre-commit run -a

test:
	go test -race -coverprofile=coverage.out -covermode=atomic

codecov: test
	codecov -t ${CODECOV_TOKEN}
