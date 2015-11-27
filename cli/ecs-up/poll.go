package main

import (
	"log"
	"time"

	"github.com/99designs/ecs-cli/cli"
)

type PollCommandInput struct {
	StackName string
}

func PollCommand(ui *cli.Ui, input PollCommandInput) {
	poller := cfnClient.PollStackEvents(input.StackName)

	log.Printf("Events for %s", input.StackName)
	for event := range poller.Events(time.Time{}) {
		ui.Printf("%s\n", event.String())
	}

	if err := poller.Error(); err != nil {
		ui.Fatal(err)
	}

	outputs, err := cfnClient.StackOutputs(input.StackName)
	if err != nil {
		ui.Fatal(err)
	}

	for k, v := range outputs {
		ui.Printf("Stack Output: %s = %s", k, v)
	}
}
