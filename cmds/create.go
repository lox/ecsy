package cmds

import (
	"os"
	"path/filepath"
	"time"

	"github.com/99designs/ecs-former/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type EcsClusterParameters struct {
	KeyName string
}

type CreateClusterCommandInput struct {
	Name       string
	StackName  string
	Templates  core.Templates
	Parameters EcsClusterParameters
}

func CreateClusterCommand(ui *Ui, input CreateClusterCommandInput) {
	if input.Name == "" {
		cwd, err := os.Getwd()
		if err != nil {
			ui.Fatal(err)
		}

		input.Name = filepath.Base(filepath.Base(cwd))
	}

	// TODO: check for existing cluster
	ui.Printf("Creating cluster %s", input.Name)
	_, err := ecs.New(session.New()).CreateCluster(&ecs.CreateClusterInput{
		ClusterName: aws.String(input.Name),
	})

	if input.StackName == "" {
		input.StackName = input.Name + "-ecs-" + time.Now().Format("20060102-150405")
	}

	ui.Printf("Creating cloudformation stack %s", input.StackName)
	svc := cloudformation.New(session.New())
	_, err = svc.CreateStack(&cloudformation.CreateStackInput{
		StackName: aws.String(input.StackName),
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
		},
		DisableRollback: aws.Bool(true),
		Parameters: []*cloudformation.Parameter{{
			ParameterKey:   aws.String("KeyName"),
			ParameterValue: aws.String(input.Parameters.KeyName),
		}, {
			ParameterKey:   aws.String("ECSCluster"),
			ParameterValue: aws.String(input.Name),
		}},
		TemplateBody: aws.String(string(input.Templates.EcsStack())),
	})
	if err != nil {
		ui.Fatal(err)
	}

	events, errs := pollStack(svc, input.StackName)
	for event := range events {
		ui.Println(eventString(event))
	}

	if err := <-errs; err != nil {
		ui.Fatal(err)
	}
}
