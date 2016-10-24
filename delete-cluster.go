package main

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/lox/ecsy/api"
)

func DeleteClusterCommand(ui *Ui, clusterName string) {
	stacks, err := api.FindAllStacksForCluster(cfnSvc, clusterName)
	if err != nil {
		ui.Fatal(err)
	}

	for _, stack := range stacks {
		err := api.DeleteStack(cfnSvc, *stack.StackName)
		if err != nil {
			ui.Fatal(err)
		}

		err = api.PollUntilDeleted(cfnSvc, *stack.StackName, func(event *cloudformation.StackEvent) {
			ui.Printf("%s\n", api.FormatStackEvent(event))
		})

		ui.Printf("Deleted %s\n", *stack.StackName)
	}
}
