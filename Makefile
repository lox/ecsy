PREFIX=github.com/99designs/ecs-cli
VERSION=$(shell git describe --tags --candidates=1 --dirty 2>/dev/null || echo "dev")
FLAGS=-X main.Version=$(VERSION)

build: templates
	go build -ldflags="$(FLAGS)" github.com/99designs/ecs-cli/cli/ecs-deploy
	go build -ldflags="$(FLAGS)" github.com/99designs/ecs-cli/cli/ecs-up

install: templates
	go install -ldflags="$(FLAGS)" github.com/99designs/ecs-cli/cli/ecs-deploy
	go install -ldflags="$(FLAGS)" github.com/99designs/ecs-cli/cli/ecs-up

templates: cloudformation/templates/build/ecs-service.json cloudformation/templates/build/ecs-stack.json cloudformation/templates/build/vpc.json
	esc -o cloudformation/templates/static.go -pkg templates cloudformation/templates/build

clean:
	rm $(wildcard cloudformation/templates/build/*.json)

cloudformation/templates/build/ecs-stack.json: cloudformation/templates/src/ecs-stack.yml cloudformation/templates/src/vpc.yml
	@mkdir -p cloudformation/templates/build/
	cfoo $^ > $@

cloudformation/templates/build/vpc.json: cloudformation/templates/src/vpc.yml
	@mkdir -p cloudformation/templates/build/
	cfoo $^ > $@

cloudformation/templates/build/ecs-service.json: cloudformation/templates/src/ecs-service.yml
	@mkdir -p cloudformation/templates/build/
	cfoo $^ > $@