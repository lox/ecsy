PREFIX=github.com/lox/ecsy
VERSION=$(shell git describe --tags --candidates=1 --dirty 2>/dev/null || echo "dev")
FLAGS=-X main.Version=$(VERSION)
TEMPLATES=templates/src/ecs-service.yml templates/src/ecs-stack.yml templates/src/network-stack.yml

.PHONY: test setup build install clean templates

test:
	govendor test +local

setup:
	go get -u github.com/kardianos/govendor
	go get -u github.com/mjibson/esc

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

release:
	equinox release \
		--version=$(VERSION) \
		--platforms='darwin_amd64 linux_amd64' \
		--signing-key=$(EQUINOX_KEY) \
		--app=app_3dH3rVQXD5e \
		--token=$(EQUINOX_TOKEN) \
		-- -ldflags='-X main.Version=$(VERSION) -s -w' \
		$(PREFIX)
