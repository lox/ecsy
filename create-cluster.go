package main

import (
	"strconv"
	"time"

	"github.com/99designs/ecs-cli/templates"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/99designs/ecs-cli/api"
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

func CreateClusterCommand(ui *Ui, input CreateClusterCommandInput) {
	ui.Printf("Creating cluster %s", input.ClusterName)

	_, err := ecsSvc.CreateCluster(&ecs.CreateClusterInput{
		ClusterName: aws.String(input.ClusterName),
	})

	network := getOrCreateNetworkStack(ui, input.ClusterName)

	timer := time.Now()
	stackName := input.ClusterName + "-ecs-" + time.Now().Format("20060102-150405")
	ui.Printf("Creating cloudformation stack %s", stackName)

	err = api.CreateStack(cfnSvc, stackName, templates.EcsStack(), map[string]string{
		"Subnets":            network.Subnets,
		"SecurityGroup":      network.SecurityGroup,
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

	err = api.PollUntilCreated(cfnSvc, stackName, func(event *cloudformation.StackEvent) {
		ui.Printf("%s\n", api.FormatStackEvent(event))
	})
	if err != nil {
		ui.Fatal(err)
	}

	ui.Printf("Cluster %s created in %s\n\n", input.ClusterName, time.Now().Sub(timer).String())
}

func getOrCreateNetworkStack(ui *Ui, clusterName string) api.NetworkOutputs {
	outputs, err := api.FindNetworkStack(cfnSvc, clusterName)

	if err != nil {
		timer := time.Now()
		ui.Printf("Creating Network Stack for %s", clusterName)

		err = api.CreateStack(cfnSvc, outputs.StackName, templates.NetworkStack(), map[string]string{})
		if err != nil {
			ui.Fatal(err)
		}

		err = api.PollUntilCreated(cfnSvc, outputs.StackName, func(event *cloudformation.StackEvent) {
			ui.Printf("%s\n", api.FormatStackEvent(event))
		})
		if err != nil {
			ui.Fatal(err)
		}

		outputs, err = api.FindNetworkStack(cfnSvc, clusterName)
		if err != nil {
			ui.Fatal(err)
		}

		ui.Printf("%s created in %s\n\n", outputs.StackName, time.Now().Sub(timer).String())
	}

	return outputs
}
