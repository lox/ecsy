ECSy
=============

A tool for managing and deploying ECS clusters, because the [official one](https://github.com/aws/amazon-ecs-cli) is [terrible](#why-not-amazon-ecs-cli)

Derived from `99designs/ecs-cli`, many thanks to those guys for being awesome and making it possible for me to release it open-source. 

## Features 

 * CloudFormation templates are used to manage the network stack (VPC), ECS services and ECS cluster
 * Designed for managing many people who manage many different ECS clusters 

## Installing

```
go get github.com/lox/ecsy
```

## Usage

### Create a new ECS Cluster with your app running

```bash
# create an ecs cluster and supporting infrastructure (vpc, autoscale group, security groups, etc)
ecsy create-cluster --cluster example --keyname lox --type m4.large --count 4

# create an ecs task and service from a docker-compose file
ecsy create-service --cluster example -f docker-compose.yml
```

### Deploy a new release of your app to a service created above

```bash
# Creates and deploys a new task with the helloworld container updated with a new image tag
ecsy deploy --cluster example -f docker-compose.yml helloworld=:v2
```

### Why not amazon-ecs-cli?

The main issue with `amazon-ecs-cli` is that it tries to emulate the `docker-compose` interface, which isn't a sensible abstraction and ends up making the architecture overly complicated. Contributing the changes we wanted upstream just wasn't viable, and beyond that issues go unanswered and development seems stagnant:

- https://github.com/aws/amazon-ecs-cli/issues/90
- https://github.com/aws/amazon-ecs-cli/issues/21
