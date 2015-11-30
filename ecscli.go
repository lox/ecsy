package ecscli

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func FindClusterStack(svc *cloudformation.CloudFormation, clusterName string) (*cloudformation.Stack, error) {
	clusterStack, err := FindStackByOutputs(svc, map[string]string{
		"StackType":  "ecs-former::ecs-stack",
		"ECSCluster": clusterName,
	})
	if err != nil {
		if err == ErrNoStacksFound {
			return nil, fmt.Errorf(
				"Failed to find a cloudformation stack for cluster %q", clusterName)
		} else {
			return nil, err
		}
	}

	if err = StackOutputMap(clusterStack).RequireKeys("Subnets", "Vpc", "SecurityGroup"); err != nil {
		return nil, err
	}

	return clusterStack, err
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
