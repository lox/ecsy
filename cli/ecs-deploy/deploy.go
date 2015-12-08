package main

import (
	"log"
	"time"

	"github.com/99designs/ecs-cli"
	"github.com/99designs/ecs-cli/cli"
	"github.com/99designs/ecs-cli/compose"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type DeployCommandInput struct {
	ClusterName     string
	ProjectName     string
	ComposeFile     string
	ContainerImages map[string]string
	HealthCheckUrl  string
}

func DeployCommand(ui *cli.Ui, input DeployCommandInput) {
	ui.Printf("Generating task definition from %s", input.ComposeFile)
	taskDefinitionInput, err := compose.TransformComposeFile(input.ComposeFile, input.ProjectName)
	if err != nil {
		ui.Fatal(err)
	}

	ui.Printf("Updating task definition for task %s", *taskDefinitionInput.Family)
	err = ecscli.UpdateContainerImages(taskDefinitionInput.ContainerDefinitions, input.ContainerImages)
	if err != nil {
		ui.Fatal(err)
	}

	resp, err := ecsSvc.RegisterTaskDefinition(taskDefinitionInput)
	if err != nil {
		ui.Fatal(err)
	}
	ui.Printf("Registered task definition %s:%d", *resp.TaskDefinition.Family, *resp.TaskDefinition.Revision)

	serviceStack, err := ecscli.FindServiceStack(cfnSvc, input.ClusterName, input.ProjectName)
	if err != nil {
		ui.Fatal(err)
	}
	log.Printf("Found service stack %s", *serviceStack.StackName)

	outputs := ecscli.StackOutputMap(serviceStack)
	timer := time.Now()

	ui.Printf("Updating service %s with new task definition", *serviceStack.StackName)
	updateResp, err := ecsSvc.UpdateService(&ecs.UpdateServiceInput{
		Service:        aws.String(outputs["ECSService"]),
		Cluster:        aws.String(outputs["ECSCluster"]),
		TaskDefinition: aws.String(*resp.TaskDefinition.TaskDefinitionArn),
	})
	if err != nil {
		ui.Fatal(err)
	}

	ui.Printf("Waiting for service to reach a steady state.")
	err = ecscli.PollDeployment(ecsSvc, updateResp.Service, outputs["ECSCluster"], func(e *ecs.ServiceEvent) {
		log.Println(*e.Message)
	})

	log.Printf("Deployed in %s", time.Now().Sub(timer).String())

	// serviceOutputs := serviceStack.OutputMap()

	// ui.Printf("Updating service %s with %s", serviceOutputs["ECSService"], taskDef.String())
	// err = ecsClient.UpdateService(serviceOutputs["ECSCluster"], serviceOutputs["ECSService"], taskDef.Arn)
	// if err != nil {
	// 	ui.Fatal(err)
	// }

	// log.Printf("%#v", resp)

	// ecscli.CloudFormation.

	// taskDef, err := ecscli.UpdateComposeTaskDefinition(input.ComposeFile, input.ProjectName, input.ContainerImages)
	// if err != nil {
	// 	ui.Fatal(err)
	// }

	// log.Printf("%#v", taskDef)

	// taskDef, err := ecsClient.UpdateComposerTaskDefinition(input.ComposeFile, input.ProjectName, input.ImageTags)
	// if err != nil {
	// 	ui.Fatal(err)
	// }

	// serviceOutputs := serviceStack.OutputMap()

	// ui.Printf("Updating service %s with %s", serviceOutputs["ECSService"], taskDef.String())
	// err = ecsClient.UpdateService(serviceOutputs["ECSCluster"], serviceOutputs["ECSService"], taskDef.Arn)
	// if err != nil {
	// 	ui.Fatal(err)
	// }

	// ui.Printf("Waiting for service to stabilize")
	// if err = ecsClient.WaitUntilServicesStable(input.ClusterName, serviceOutputs["ECSService"]); err != nil {
	// 	ui.Fatal(err)
	// }

	// ui.Println("Service available at", serviceOutputs["ECSLoadBalancer"])

}
