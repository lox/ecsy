ECS CLI Tools
=============

This collection of tools provides a stop-gap whilst the official AWS ECS CLI tools develop.

## Installing

```
export GO15VENDOREXPERIMENT=1
go install -v -o ecs-cli github.com/99designs/ecs-cli/cli
```

## Usage

### Create a new ECS Cluster with your app running

```bash
# create an ecs cluster and supporting infrastructure (vpc, autoscale group, security groups, etc)
ecs-cli create-cluster --cluster example --keyname lox --type m4.large --count 4

# create an ecs task and service from a docker-compose file
ecs-cli create-service --cluster example -f docker-compose.yml
```

### Deploy a new release of your app to a service created above

```bash
# Creates and deploys a new task with the helloworld container updated with a new image tag
ecs-cli deploy --cluster example -f docker-compose.yml helloworld=:v2
```

### TODO

```
ecs-cli update-cluster --cluster <ecs_cluster>
ecs-cli create-service -f <docker-compose.yml> --cluster <ecs_cluster>
ecs-cli update-service -f <docker-compose.yml> --cluster <ecs_cluster>
ecs-cli run-task -f <docker-compose.yml> --cluster <ecs_cluster> <container>
```


