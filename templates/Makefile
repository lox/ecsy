
all: build/ecs-app.json build/ecs-cluster.json build/vpc.json build/ecs-cluster-byo-vpc.json

clean:
	rm $(wildcard build/*.json)

build/ecs-cluster.json: src/ecs-cluster.yml src/vpc.yml
	@mkdir -p build/
	cfoo $^ > $@

build/ecs-cluster-byo-vpc.json: src/ecs-cluster.yml src/vpc-byo.yml
	@mkdir -p build/
	cfoo $^ > $@

build/vpc.json: src/vpc.yml
	@mkdir -p build/
	cfoo $^ > $@

build/ecs-app.json: src/ecs-app.yml
	@mkdir -p build/
	cfoo $^ > $@