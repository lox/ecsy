ECS CLI Tools
=============

This collection of tools provides a stop-gap whilst the official AWS ECS CLI tools develop.

They use cloudformations internally, which can be used without the tooling if required.

## Installing



## Usage

### Create a new ECS Cluster with your app running

```bash
# create an ecs cluster and supporting infrastructure (vpc, autoscale group, security groups, etc)
ecs-up create-cluster --cluster example --keyname lox --type m4.large --size 4

# create an ecs task and service from a docker-compose file
ecs-up create-service --cluster example -f docker-compose.yml
```

### Deploy a new release of your app to a service created above

```bash
# Creates and deploys a new task with the helloworld container updated with a new image tag
ecs-deploy --cluster example -f docker-compose.yml helloworld=:v2
```

### TODO

```
ecs-up update-cluster --cluster <ecs_cluster>
ecs-up create-service -f <docker-compose.yml> --cluster <ecs_cluster>
ecs-up update-service -f <docker-compose.yml> --cluster <ecs_cluster>
ecs-docker-run -f <docker-compose.yml> --cluster <ecs_cluster> <container>
```


