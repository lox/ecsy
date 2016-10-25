PREFIX=github.com/lox/ecsy
VERSION=$(shell git describe --tags --candidates=1 --dirty 2>/dev/null || echo "dev")
FLAGS=-X main.Version=$(VERSION)

test:
	go get github.com/kardianos/govendor
	govendor test +local

setup:
	gem install cfoo
	go get github.com/mjibson/esc

build: templates
	go build -ldflags="$(FLAGS)" $(PREFIX)

install: templates
	go install -ldflags="$(FLAGS)" $(PREFIX)

templates: templates/build/ecs-service.json templates/build/ecs-stack.json templates/build/network-stack.json
	esc -o templates/static.go -pkg templates templates/build

clean:
	rm $(wildcard templates/build/*.json)

templates/build/ecs-stack.json: templates/src/ecs-stack.yml
	@mkdir -p templates/build/
	cfoo $^ > $@
	test -s $@

templates/build/network-stack.json: templates/src/network-stack.yml
	@mkdir -p templates/build/
	cfoo $^ > $@
	test -s $@

templates/build/ecs-service.json: templates/src/ecs-service.yml
	@mkdir -p templates/build/
	cfoo $^ > $@
	test -s $@