//go:generate esc -o core/static.go -pkg core templates/build
package main

import (
	"io/ioutil"
	"log"

	"github.com/99designs/ecs-former/cmds"
	"github.com/99designs/ecs-former/core"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	Version string
)

func main() {
	var (
		debug         = kingpin.Flag("debug", "Show debugging output").Bool()
		create        = kingpin.Command("create-cluster", "Create an ECS cluster")
		createName    = create.Flag("name", "The name of the ECS stack to create").String()
		createKeyName = create.Flag("keyname", "The EC2 keypair to use for instance").Default("default").String()
		poll          = kingpin.Command("poll", "Poll a CloudFormation stack")
		pollStack     = poll.Arg("stack", "The name of the stack to poll").Required().String()
	)

	kingpin.Version(Version)
	kingpin.CommandLine.Help =
		`Manage Amazon's ECS with CloudFormation`

	ui := cmds.NewUi()
	cmd := kingpin.Parse()

	if *debug {
		log.SetFlags(0)
		log.SetOutput(ui.DebugWriter())
	} else {
		log.SetOutput(ioutil.Discard)
	}

	switch cmd {
	case create.FullCommand():
		cmds.CreateClusterCommand(ui, cmds.CreateClusterCommandInput{
			Name:      *createName,
			Templates: core.Templates{core.FS(false)},
			Parameters: cmds.EcsClusterParameters{
				KeyName: *createKeyName,
			},
		})
	case poll.FullCommand():
		cmds.PollCommand(ui, cmds.PollCommandInput{
			StackName: *pollStack,
		})
	}
}
