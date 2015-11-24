
all: build/ecs-service.json build/ecs-stack.json build/vpc.json

clean:
	rm $(wildcard build/*.json)

build/ecs-stack.json: src/ecs-stack.yml src/vpc.yml
	@mkdir -p build/
	cfoo $^ > $@

build/vpc.json: src/vpc.yml
	@mkdir -p build/
	cfoo $^ > $@

build/ecs-service.json: src/ecs-service.yml src/ecs-stack.yml src/vpc.yml
	@mkdir -p build/
	cfoo $^ > $@