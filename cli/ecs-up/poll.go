package main

import (
	"github.com/99designs/ecs-cli"
	"github.com/99designs/ecs-cli/cli"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type PollCommandInput struct {
	StackName string
}

func PollCommand(ui *cli.Ui, input PollCommandInput) {
	err := ecscli.PollStackEvents(cfnSvc, input.StackName, func(event *cloudformation.StackEvent) {
		ui.Printf("%s\n", ecscli.FormatStackEvent(event))
	})
	if err != nil {
		ui.Fatal(err)
	}

	outputs, err := ecscli.StackOutputs(cfnSvc, input.StackName)
	if err != nil {
		ui.Fatal(err)
	}

	for k, v := range outputs {
		ui.Printf("Stack Output: %s = %s", k, v)
	}
}
