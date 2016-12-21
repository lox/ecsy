#!/bin/bash
set -euo pipefail

make setup build

# Make sure nothing is around before we start.
echo "~~~ :aws: Delete any existing clusters"
./ecsy delete-cluster --cluster=ecsy-test

echo "--- :aws: Create cluster"
./ecsy create-cluster --cluster=ecsy-test --count=2 --type t2.nano --keyname "${KEYNAME:-default}"

echo "--- :aws: Create service"
./ecsy create-service --cluster=ecsy-test -p app -f ./examples/helloworld/docker-compose.yml

echo "--- :aws: Running once-off task"
./ecsy run-task --cluster=ecsy-test -f examples/helloworld/docker-compose.yml -s sample -- echo hello world

echo "--- :aws: Deploy update"
./ecsy deploy --cluster=ecsy-test -p app -f ./examples/helloworld/docker-compose.yml

echo "--- :aws: Delete cluster"
./ecsy delete-cluster --cluster=ecsy-test