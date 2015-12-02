package main

import (
	"log"
	"strconv"
	"time"

	"github.com/99designs/ecs-cli"
	"github.com/99designs/ecs-cli/cli"
	"github.com/99designs/ecs-cli/templates"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type EcsClusterParameters struct {
	KeyName            string
	InstanceType       string
	InstanceCount      int
	DockerHubUsername  string
	DockerHubPassword  string
	DockerHubEmail     string
	DatadogApiKey      string
	LogspoutTarget     string
	AuthorizedUsersUrl string
}

type CreateClusterCommandInput struct {
	ClusterName string
	Parameters  EcsClusterParameters
}

func CreateClusterCommand(ui *cli.Ui, input CreateClusterCommandInput) {
	ui.Printf("Creating cluster %s", input.ClusterName)

	_, err := ecsSvc.CreateCluster(&ecs.CreateClusterInput{
		ClusterName: aws.String(input.ClusterName),
	})

	timer := time.Now()
	stackName := input.ClusterName + "-ecs-" + time.Now().Format("20060102-150405")
	ui.Printf("Creating cloudformation stack %s", stackName)

	err = ecscli.CreateStack(cfnSvc, stackName, templates.EcsStack(), map[string]string{
		"KeyName":            input.Parameters.KeyName,
		"ECSCluster":         input.ClusterName,
		"InstanceType":       input.Parameters.InstanceType,
		"DesiredCapacity":    strconv.Itoa(input.Parameters.InstanceCount),
		"DockerHubUsername":  input.Parameters.DockerHubUsername,
		"DockerHubPassword":  input.Parameters.DockerHubPassword,
		"DockerHubEmail":     input.Parameters.DockerHubEmail,
		"LogspoutTarget":     input.Parameters.LogspoutTarget,
		"DatadogApiKey":      input.Parameters.DatadogApiKey,
		"AuthorizedUsersUrl": input.Parameters.AuthorizedUsersUrl,
	})
	if err != nil {
		ui.Fatal(err)
	}

	err = ecscli.PollStackEvents(cfnSvc, stackName, func(event *cloudformation.StackEvent) {
		ui.Printf("%s\n", ecscli.FormatStackEvent(event))
	})
	if err != nil {
		ui.Fatal(err)
	}

	log.Printf("Cluster %s created in %s", input.ClusterName, time.Now().Sub(timer).String())
}
