package main

import (
	"log"
	"time"

	"github.com/99designs/ecs-cli/cli"
	"github.com/99designs/ecs-cli/cloudformation"
	"github.com/99designs/ecs-cli/cloudformation/templates"
)

type CreateClusterCommandInput struct {
	ClusterName string
	Parameters  EcsClusterParameters
}

func CreateClusterCommand(ui *cli.Ui, input CreateClusterCommandInput) {
	ui.Printf("Creating cluster %s", input.ClusterName)

	_, err := ecsClient.CreateCluster(input.ClusterName)
	if err != nil {
		ui.Fatal(err)
	}

	timer := time.Now()
	stackName := input.ClusterName + "-ecs-" + time.Now().Format("20060102-150405")
	ui.Printf("Creating cloudformation stack %s", stackName)

	err = cfnClient.CreateStack(stackName, templates.EcsStack(), cloudformation.StackParameters{
		"KeyName":    input.Parameters.KeyName,
		"ECSCluster": input.ClusterName,
	})
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

	outputs, err := cfnClient.StackOutputs(stackName)
	if err != nil {
		ui.Fatal(err)
	}

	for k, v := range outputs {
		log.Printf("Stack Output: %s = %s", k, v)
	}

	log.Printf("Cluster %s created in %s", input.ClusterName, time.Now().Sub(timer).String())
}
