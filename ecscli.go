package ecscli

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type NetworkOutputs struct {
	StackName     string
	Vpc           string
	Subnets       string
	SecurityGroup string
}

func FindServiceStack(svc *cloudformation.CloudFormation, clusterName, taskFamily string) (*cloudformation.Stack, error) {
	serviceStack, err := FindStackByOutputs(svc, map[string]string{
		"StackType":  "ecs-former::ecs-service",
		"ECSCluster": clusterName,
		"TaskFamily": taskFamily,
	})
	if err != nil {
		if err == ErrNoStacksFound {
			return nil, fmt.Errorf(
				"Failed to find a cloudformation stack for task %q, cluster %q", taskFamily, clusterName)
		}
	}
	return serviceStack, err
}

func FindNetworkStack(svc cfnInterface, clusterName string) (NetworkOutputs, error) {
	stackName := clusterName + "-network"

	outputs, err := StackOutputs(svc, stackName)
	if err != nil {
		return NetworkOutputs{StackName: stackName}, err
	}

	if err := outputs.RequireKeys("Vpc", "Subnets", "SecurityGroup"); err != nil {
		return NetworkOutputs{StackName: stackName}, err
	}

	return NetworkOutputs{
		StackName:     stackName,
		Vpc:           outputs["Vpc"],
		Subnets:       outputs["Subnets"],
		SecurityGroup: outputs["SecurityGroup"],
	}, nil
}
