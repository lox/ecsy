package cmd

import (
	"log"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/lox/ecsy/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigurePollStack(app *kingpin.Application, svc api.Services) {
	var stackName string

	cmd := app.Command("poll-stack", "Poll a cloudformation stack until it's complete")
	cmd.Alias("poll")

	cmd.Flag("stack", "The name of the cloudformation stack to poll").
		Required().
		StringVar(&stackName)

	cmd.Action(func(c *kingpin.ParseContext) error {
		err := api.PollUntilCreated(svc.Cloudformation, stackName, func(event *cloudformation.StackEvent) {
			log.Printf("%s\n", api.FormatStackEvent(event))
		})
		if err != nil {
			return err
		}

		outputs, err := api.StackOutputs(svc.Cloudformation, stackName)
		if err != nil {
			return err
		}

		for k, v := range outputs {
			log.Printf("Stack Output: %s = %s", k, v)
		}

		return nil
	})
}
