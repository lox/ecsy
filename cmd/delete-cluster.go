package cmd

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/lox/ecsy/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureDeleteCluster(app *kingpin.Application, svc api.Services) {
	var cluster string

	cmd := app.Command("delete-cluster", "Deletes a cluster and all running services on it")
	cmd.Flag("cluster", "The name of the ECS cluster to delete").
		Required().
		StringVar(&cluster)

	cmd.Action(func(c *kingpin.ParseContext) error {
		fmt.Printf("Deleting cluster %s", cluster)
		timer := time.Now()

		stacks, err := api.FindAllStacksForCluster(svc.Cloudformation, cluster)
		if err != nil {
			return err
		}

		for _, stack := range stacks {
			err := api.DeleteStack(svc.Cloudformation, *stack.StackName)
			if err != nil {
				return err
			}

			err = api.PollUntilDeleted(svc.Cloudformation, *stack.StackName, func(event *cloudformation.StackEvent) {
				fmt.Printf("%s\n", api.FormatStackEvent(event))
			})

			fmt.Printf("Deleted stack %s\n", *stack.StackName)
		}

		fmt.Printf("Cluster %s deleted in %s\n\n", cluster, time.Now().Sub(timer).String())
		return nil
	})
}
