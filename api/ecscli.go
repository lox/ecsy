package api

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type NetworkOutputs struct {
	StackName      string
	VpcId          string
	Subnet0Public  string
	Subnet1Public  string
	Subnet2Private string
	Subnet3Private string
}

func FindClusterStack(svc cfnInterface, clusterName string) (*cloudformation.Stack, error) {
	clusterStacks, err := FindStacksByOutputs(svc, map[string]string{
		"StackType":  "ecs-former::ecs-stack",
		"ECSCluster": clusterName,
	})
	if err != nil {
		return nil, err
	}
	if len(clusterStacks) == 0 {
		return nil, fmt.Errorf(
			"Failed to find a cloudformation stack for cluster %q",
			clusterName,
		)
	}
	return clusterStacks[0], nil
}

func FindServiceStack(svc cfnInterface, clusterName, taskFamily string) (*cloudformation.Stack, error) {
	serviceStacks, err := FindStacksByOutputs(svc, map[string]string{
		"StackType":  "ecs-former::ecs-service",
		"ECSCluster": clusterName,
		"TaskFamily": taskFamily,
	})
	if err != nil {
		return nil, err
	}
	if len(serviceStacks) == 0 {
		return nil, fmt.Errorf(
			"Failed to find a cloudformation stack matching task %q, cluster %q",
			taskFamily,
			clusterName,
		)
	}
	return serviceStacks[0], nil
}

func FindNetworkStack(svc cfnInterface, clusterName string) (NetworkOutputs, error) {
	stackName := clusterName + "-network"

	outputs, err := StackOutputs(svc, stackName)
	if err != nil {
		return NetworkOutputs{StackName: stackName}, err
	}

	return NetworkOutputs{
		StackName:      stackName,
		VpcId:          outputs["VpcId"],
		Subnet0Public:  outputs["Subnet0Public"],
		Subnet1Public:  outputs["Subnet1Public"],
		Subnet2Private: outputs["Subnet2Private"],
		Subnet3Private: outputs["Subnet3Private"],
	}, nil
}

func FindAllStacksForCluster(svc cfnInterface, clusterName string) ([]*cloudformation.Stack, error) {
	stacks, err := FindStacksByOutputs(svc, map[string]string{
		"ECSCluster": clusterName,
	})

	networkStack, err := FindStacksByName(svc, clusterName+"-network")
	if err == nil {
		stacks = append(stacks, networkStack...)
	}

	return stacks, nil
}
