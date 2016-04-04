#!/bin/bash

set -eu

make build

# Make sure nothing is around before we start.
echo "--- Delete cluster"
./ecs-cli delete-cluster --cluster=ecs-cli-buildkite

echo "--- Create cluster"
./ecs-cli create-cluster --cluster=ecs-cli-buildkite --count=2

echo "--- Create service"
./ecs-cli create-service --cluster=ecs-cli-buildkite -p app -f ./examples/helloworld/docker-compose.yml

echo "--- Deploy update"
./ecs-cli deploy --cluster=ecs-cli-buildkite -p app -f ./examples/helloworld/docker-compose.yml

echo "--- Delete cluster"
./ecs-cli delete-cluster --cluster=ecs-cli-buildkite
