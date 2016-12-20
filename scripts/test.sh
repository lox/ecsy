#!/bin/bash
set -euo pipefail

make setup build

# Make sure nothing is around before we start.
echo "--- Delete any existing clusters"
./ecsy delete-cluster --cluster=ecsy-test

echo "--- Create cluster"
./ecsy create-cluster --cluster=ecsy-test --count=2 --type t2.nano --keyname "${KEYNAME:-default}"

echo "--- Create service"
./ecsy create-service --cluster=ecsy-test -p app -f ./examples/helloworld/docker-compose.yml

echo "--- Deploy update"
./ecsy deploy --cluster=ecsy-test -p app -f ./examples/helloworld/docker-compose.yml

echo "--- Delete cluster"
./ecsy delete-cluster --cluster=ecsy-test