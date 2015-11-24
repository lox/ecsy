ECS Former
==========

A golang binary wrapper around ECS + CloudFormation. Creates a `ecs.json` file in the project directory to track context of what stack and app to target.

## Usage

```bash
# create VPC, ECS Cluster, EC2 AutoScale group
ecs-former up --keyname lox
```

## Plans

### `ecs-former config`

Manipulate the `ecs.json` file, or populate it from a given stack

### `ecs-former scale <number> <instance type>`

Scale out the underlying instance pool

### `ecs-former deploy`

Deploy the current revision of the with a blue-green strategy.

### `ecs-former exec`

Executes a command in a single container synchronously.



