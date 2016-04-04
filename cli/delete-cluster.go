package main

import (
	"github.com/99designs/ecs-cli"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func DeleteClusterCommand(ui *Ui, clusterName string) {
	stacks, err := ecscli.FindAllStacksForCluster(cfnSvc, clusterName)
	if err != nil {
		ui.Fatal(err)
	}

	for _, stack := range stacks {
		err := ecscli.DeleteStack(cfnSvc, *stack.StackName)
		if err != nil {
			ui.Fatal(err)
		}

		err = ecscli.PollUntilDeleted(cfnSvc, *stack.StackName, func(event *cloudformation.StackEvent) {
			ui.Printf("%s\n", ecscli.FormatStackEvent(event))
		})

		ui.Printf("Deleted %s\n", *stack.StackName)
	}
}
