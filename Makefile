PREFIX=github.com/lox/ecsy
VERSION=$(shell git describe --tags --candidates=1 --dirty 2>/dev/null || echo "dev")
FLAGS=-X main.Version=$(VERSION)
TEMPLATES=templates/src/ecs-service.yml templates/src/ecs-stack.yml templates/src/network-stack.yml

.PHONY: test setup build install clean templates

test:
	govendor test +local

setup:
	go get github.com/kardianos/govendor
	go get github.com/mjibson/esc

build: templates
	go build -ldflags="$(FLAGS)" $(PREFIX)

install: templates
	go install -ldflags="$(FLAGS)" $(PREFIX)

templates: $(TEMPLATES)
	esc -o templates/static.go -pkg templates templates/src

validate:
	@echo $(TEMPLATES) | xargs -n1 -t -I{} aws cloudformation validate-template --template-body file://{}

clean:
	rm templates/static.go