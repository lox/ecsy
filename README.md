ECS Former
==========

## Usage

```bash
ecs-former create-cluster --name wordpress-ecs
ecs-former deploy --cluster wordpress-ecs --taskfile examples/wordpress.task.json latest
```

## Plans

### `ecs-former scale <number> <instance type>`

Scale out the underlying instance pool

### `ecs-former exec`

Executes a command in a single container synchronously.