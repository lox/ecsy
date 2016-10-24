package main

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/lox/ecsy/api"
)

type PollCommandInput struct {
	StackName string
}

func PollCommand(ui *Ui, input PollCommandInput) {
	err := api.PollUntilCreated(cfnSvc, input.StackName, func(event *cloudformation.StackEvent) {
		ui.Printf("%s\n", api.FormatStackEvent(event))
	})
	if err != nil {
		ui.Fatal(err)
	}

	outputs, err := api.StackOutputs(cfnSvc, input.StackName)
	if err != nil {
		ui.Fatal(err)
	}

	for k, v := range outputs {
		ui.Printf("Stack Output: %s = %s", k, v)
	}
}
