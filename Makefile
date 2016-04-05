PREFIX=github.com/99designs/ecs-cli
VERSION=$(shell git describe --tags --candidates=1 --dirty 2>/dev/null || echo "dev")
FLAGS=-X main.Version=$(VERSION)

vendor: glide.lock glide.yaml
	glide -q install
	touch vendor

build: templates vendor
	go build -ldflags="$(FLAGS)" github.com/99designs/ecs-cli

install: templates
	go install -ldflags="$(FLAGS)" github.com/99designs/ecs-cli

templates: templates/build/ecs-service.json templates/build/ecs-stack.json templates/build/network-stack.json
	esc -o templates/static.go -pkg templates templates/build

clean:
	rm $(wildcard templates/build/*.json)

templates/build/ecs-stack.json: templates/src/ecs-stack.yml
	@mkdir -p templates/build/
	cfoo $^ > $@

templates/build/network-stack.json: templates/src/network-stack.yml
	@mkdir -p templates/build/
	cfoo $^ > $@

templates/build/ecs-service.json: templates/src/ecs-service.yml
	@mkdir -p templates/build/
	cfoo $^ > $@
