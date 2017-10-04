#!/bin/bash
set -euo pipefail

make setup build

# Make sure nothing is around before we start.
echo "~~~ :aws: :ecs: Delete any existing clusters"
./ecsy delete-cluster --cluster=ecsy-test

echo "--- :aws: :ecs: Create cluster"
./ecsy create-cluster --cluster=ecsy-test --count=2 --type t2.nano --keyname "${KEYNAME:-default}"

echo "--- :aws: :ecs: Create service"
./ecsy create-service --cluster=ecsy-test --name app -f ./examples/helloworld/taskdefinition.json

echo "--- :aws: :ecs: Running once-off task"
./ecsy run-task --cluster=ecsy-test -f ./examples/helloworld/taskdefinition.json --service sample -- echo hello world

echo "--- :aws: :ecs: Deploy update"
./ecsy deploy --cluster=ecsy-test --name app -f ./examples/helloworld/taskdefinition.json

echo "--- :aws: :ecs: Delete cluster"
./ecsy delete-cluster --cluster=ecsy-test
