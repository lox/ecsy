package cmds

import (
	"log"
	"path/filepath"
	"time"

	"github.com/99designs/ecs-former/core"
	"github.com/99designs/ecs-former/transform"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type DeployCommandInput struct {
	ComposeFile   string
	ProjectName   string
	ClusterName   string
	ContainerName string
	Templates     core.Templates
}

func DeployCommand(ui *Ui, input DeployCommandInput) {
	if input.ProjectName == "" {
		input.ProjectName = filepath.Base(filepath.Dir(input.ComposeFile))
	}

	cfnSvc := cloudformation.New(session.New())
	clusterStack, err := findStack(cfnSvc, map[string]string{
		"StackType":  "ecs-former::ecs-stack",
		"ECSCluster": input.ClusterName,
	})
	if err != nil {
		ui.Fatal(err)
	}
	if clusterStack == nil {
		ui.Fatalf("Failed to find stack for cluster %s", input.ClusterName)
	}
	log.Printf("Found cluster stack %s", clusterStack.StackName)

	subnets, ok := outputValue("Subnets", clusterStack.Outputs)
	if !ok {
		log.Fatal("No Subnets in stack outputs for " + *clusterStack.StackName)
	}

	vpc, ok := outputValue("Vpc", clusterStack.Outputs)
	if !ok {
		log.Fatal("No Vpc in stack outputs for " + *clusterStack.StackName)
	}

	securityGroup, ok := outputValue("SecurityGroup", clusterStack.Outputs)
	if !ok {
		log.Fatal("No SecurityGroup in stack outputs for " + *clusterStack.StackName)
	}

	regTaskInput, err := transform.TransformComposeFile(input.ComposeFile, input.ProjectName)
	if err != nil {
		ui.Fatal(err)
	}

	log.Printf("%#v", regTaskInput)
	svc := ecs.New(session.New(nil))

	ui.Println("Registering task definition with ECS")
	resp, err := svc.RegisterTaskDefinition(regTaskInput)
	if err != nil {
		ui.Fatal(err)
	}

	log.Printf("Registered task definition %s:%d (%s)",
		*resp.TaskDefinition.Family,
		*resp.TaskDefinition.Revision,
		*resp.TaskDefinition.TaskDefinitionArn)

	// look for the service stack
	serviceStack, err := findStack(cfnSvc, map[string]string{
		"StackType":  "ecs-former::ecs-service",
		"ECSCluster": input.ClusterName,
		"TaskFamily": *resp.TaskDefinition.Family,
	})
	if err != nil {
		ui.Fatal(err)
	}

	var stackName string

	if serviceStack == nil {
		stackName = input.ProjectName + "-ecs-service-" + time.Now().Format("20060102-150405")
		ui.Printf("Failed to find an existing service stack, creating %s", stackName)

		_, err = cfnSvc.CreateStack(&cloudformation.CreateStackInput{
			StackName: aws.String(stackName),
			Capabilities: []*string{
				aws.String("CAPABILITY_IAM"),
			},
			Parameters: []*cloudformation.Parameter{{
				ParameterKey:   aws.String("ECSCluster"),
				ParameterValue: aws.String(input.ClusterName),
			}, {
				ParameterKey:   aws.String("TaskFamily"),
				ParameterValue: resp.TaskDefinition.Family,
			}, {
				ParameterKey:   aws.String("TaskDefinition"),
				ParameterValue: resp.TaskDefinition.TaskDefinitionArn,
			}, {
				ParameterKey:   aws.String("ContainerName"),
				ParameterValue: aws.String(input.ContainerName),
			}, {
				ParameterKey:   aws.String("Subnets"),
				ParameterValue: aws.String(subnets),
			}, {
				ParameterKey:   aws.String("Vpc"),
				ParameterValue: aws.String(vpc),
			}, {
				ParameterKey:   aws.String("ECSSecurityGroup"),
				ParameterValue: aws.String(securityGroup),
			}},
			TemplateBody: aws.String(string(input.Templates.EcsService())),
		})
		if err != nil {
			ui.Fatal(err)
		}
	} else {
		stackName := *serviceStack.StackName
		ui.Printf("Found existing service stack, updating %s", stackName)

		_, err := cfnSvc.UpdateStack(&cloudformation.UpdateStackInput{
			StackName: aws.String(stackName),
			Capabilities: []*string{
				aws.String("CAPABILITY_IAM"),
			},
			Parameters: []*cloudformation.Parameter{{
				ParameterKey:   aws.String("ECSCluster"),
				ParameterValue: aws.String(input.ClusterName),
			}, {
				ParameterKey:   aws.String("TaskFamily"),
				ParameterValue: resp.TaskDefinition.Family,
			}, {
				ParameterKey:   aws.String("TaskDefinition"),
				ParameterValue: resp.TaskDefinition.TaskDefinitionArn,
			}, {
				ParameterKey:   aws.String("ContainerName"),
				ParameterValue: aws.String(input.ContainerName),
			}, {
				ParameterKey:   aws.String("Subnets"),
				ParameterValue: aws.String(subnets),
			}, {
				ParameterKey:   aws.String("Vpc"),
				ParameterValue: aws.String(vpc),
			}, {
				ParameterKey:   aws.String("ECSSecurityGroup"),
				ParameterValue: aws.String(securityGroup),
			}},
			UsePreviousTemplate: aws.Bool(true),
		})
		if err != nil {
			ui.Fatal(err)
		}
	}

	events, errs := pollStack(cfnSvc, stackName)
	for event := range events {
		ui.Println(eventString(event))
	}

	if err := <-errs; err != nil {
		ui.Fatal(err)
	}

}

func outputValue(key string, outputs []*cloudformation.Output) (val string, exist bool) {
	for _, output := range outputs {
		if *output.OutputKey == key {
			return *output.OutputValue, true
		}
	}
	return "", false
}

func matchStackOutputs(s *cloudformation.Stack, match map[string]string) bool {
	for k, v := range match {
		val, ok := outputValue(k, s.Outputs)
		if !ok || val != v {
			return false
		}
	}
	return true
}

func findStack(svc *cloudformation.CloudFormation, outputs map[string]string) (*cloudformation.Stack, error) {
	stacks, err := allActiveStacks(svc)
	if err != nil {
		return nil, err
	}
	for _, stack := range stacks {
		if matchStackOutputs(stack, outputs) {
			return stack, nil
		}
	}
	return nil, nil
}

func allActiveStacks(svc *cloudformation.CloudFormation) ([]*cloudformation.Stack, error) {
	stacks := []*cloudformation.Stack{}
	return stacks, svc.DescribeStacksPages(nil, func(page *cloudformation.DescribeStacksOutput, lastPage bool) bool {
		for _, s := range page.Stacks {
			if *s.StackStatus == "CREATE_COMPLETE" || *s.StackStatus == "UPDATE_COMPLETE" {
				stacks = append(stacks, s)
			}
		}
		return lastPage
	})
}
