package main

import (
	"log"
	"strconv"
	"time"

	"github.com/99designs/ecs-cli/cli"
	"github.com/99designs/ecs-cli/cloudformation"
	"github.com/99designs/ecs-cli/cloudformation/templates"
)

type CreateServiceCommandInput struct {
	ClusterName       string
	ProjectName       string
	DockerComposeFile string
}

func CreateServiceCommand(ui *cli.Ui, input CreateServiceCommandInput) {
	exists, err := isServiceCreated(input.ProjectName, input.ClusterName)
	if err != nil {
		ui.Fatal(err)
	} else if exists {
		ui.Fatalf("A service already exists for %q in cluster %q. Use `ecs-deploy` or `ecs-up update-service`",
			input.ProjectName, input.ClusterName)
	}

	ui.Printf("Registering a task for %s (%s)", input.ProjectName, input.DockerComposeFile)

	taskDef, err := ecsClient.RegisterComposerTaskDefinition(input.DockerComposeFile, input.ProjectName)
	if err != nil {
		ui.Fatal(err)
	}

	ui.Printf("Created task %s revision %d", taskDef.Family, taskDef.Revision)

	// find the ecs stack for parameters
	stack, err := cfnClient.FindStackByOutputs(cloudformation.OutputMap{
		"StackType":  "ecs-former::ecs-stack",
		"ECSCluster": input.ClusterName,
	})
	if err != nil {
		ui.Fatal(err)
	}

	log.Printf("Found cloudformation stack %s for ECS cluster", stack.StackName)

	outputs := stack.OutputMap()

	if err = outputs.RequireKeys("Subnets", "Vpc", "SecurityGroup"); err != nil {
		ui.Fatal(err)
	}

	params := cloudformation.StackParameters{
		"ECSCluster":       input.ClusterName,
		"TaskFamily":       taskDef.Family,
		"TaskDefinition":   taskDef.Arn,
		"Subnets":          outputs["Subnets"],
		"Vpc":              outputs["Vpc"],
		"ECSSecurityGroup": outputs["SecurityGroup"],
	}

	if len(taskDef.PortMappings) == 1 {
		params["ContainerName"] = taskDef.PortMappings[0].ContainerName
		params["ContainerPort"] = strconv.FormatInt(taskDef.PortMappings[0].ContainerPort, 10)
		params["ELBPort"] = strconv.FormatInt(taskDef.PortMappings[0].HostPort, 10)
	} else if len(taskDef.PortMappings) > 1 {
		ui.Fatalf("Task definition has >1 host mapped ports, not yet supported")
	}

	timer := time.Now()
	stackName := input.ClusterName + "-ecs-service-" + time.Now().Format("20060102-150405")

	ui.Printf("Creating service cloudformation stack %s", stackName)

	err = cfnClient.CreateStack(stackName, templates.EcsService(), params)
	if err != nil {
		ui.Fatal(err)
	}

	poller := cfnClient.PollStackEvents(stackName)
	for event := range poller.Events(time.Time{}) {
		ui.Printf("%s\n", event.String())
	}

	if err := poller.Error(); err != nil {
		ui.Fatal(err)
	}

	serviceOutputs, err := cfnClient.StackOutputs(stackName)
	if err != nil {
		ui.Fatal(err)
	}

	log.Printf("%#v", serviceOutputs)
	log.Printf("Service %s created in %s", input.ClusterName, time.Now().Sub(timer).String())

	log.Printf("Waiting for service to stabilize")
	if err = ecsClient.WaitUntilServicesStable(input.ClusterName, serviceOutputs["ECSService"]); err != nil {
		ui.Fatal(err)
	}

	ui.Println("Service available at", serviceOutputs["ECSLoadBalancer"])

}

func isServiceCreated(project string, cluster string) (bool, error) {
	_, err := cfnClient.FindStackByOutputs(cloudformation.OutputMap{
		"StackType":  "ecs-former::ecs-service",
		"ECSCluster": cluster,
		"TaskFamily": project,
	})
	if err != nil {
		if err == cloudformation.ErrNoStacksFound {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
