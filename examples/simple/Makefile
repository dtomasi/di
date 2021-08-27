install:
	## Install stringer
	go get -u -a golang.org/x/tools/cmd/stringer

build:
	go generate
	go build -o ./bin/greeter_app main.go

run: build
	./bin/greeter_app
