# ECSy [![Build Status](https://travis-ci.org/lox/ecsy.svg?branch=master)](https://travis-ci.org/lox/ecsy)

A tool for creating and managing ECS clusters using CloudFormation.

Originally `99designs/ecs-cli`, many thanks to those guys for being awesome and making it possible for me to release it open-source.

## Features

* CloudFormation based - Network stack, ECS cluster and ECS services
* Support for managing ECS services with ALB loadbalancers
* Designed for managing many ECS clusters
* Built-in support for common third-party services like Datadog

## Installing

Either download the binary from https://dl.equinox.io/lox/ecsy/stable, or install with golang:

```bash
go get github.com/lox/ecsy
```

## Usage

### Create a new ECS Cluster with your app running

```bash
# create an ecs cluster and supporting infrastructure (vpc, autoscale group, security groups, etc)
ecsy create-cluster --cluster example --keyname lox --type m4.large --count 4

# create a service from the provided task definition
ecsy create-service --cluster example --name example-service -f taskdefinition.json
```

### Deploy a new release of your app to a service created above

```bash
# Creates and deploys a new task with the helloworld container updated with a new image tag
ecsy deploy --cluster example --service example-service -f taskdefinition.json "helloworld=:v2"
```

## Building

Setup the build dependencies.

```
make setup
```

Install ecsy.

```
make install
```

### Why not amazon-ecs-cli?

The main issue with `amazon-ecs-cli` is that it tries to emulate the `docker-compose` interface, which isn't a sensible abstraction and ends up making the architecture overly complicated. Contributing the changes we wanted upstream just wasn't viable, and beyond that issues go unanswered and development seems stagnant:

- https://github.com/aws/amazon-ecs-cli/issues/90
- https://github.com/aws/amazon-ecs-cli/issues/21
