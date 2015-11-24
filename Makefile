PREFIX=github.com/99designs/ecs-former
VERSION=$(shell git describe --tags --candidates=1 --dirty)
FLAGS=-X main.Version=$(VERSION)

build:
	go generate
	go build -o ecs-former -ldflags="$(FLAGS)" $(PREFIX)

install:
	go generate
	go install -ldflags="$(FLAGS)" $(PREFIX)
