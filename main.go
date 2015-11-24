//go:generate esc -o core/static.go -pkg core templates/build
package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/99designs/ecs-former/cmds"
	"github.com/99designs/ecs-former/core"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	Version string
)

func main() {
	var (
		debug     = kingpin.Flag("debug", "Show debugging output").Bool()
		up        = kingpin.Command("up", "Bring your ecs cluster to a running state")
		upStack   = up.Flag("stack", "An existing cloudformation stack to use").String()
		upKeyName = up.Flag("keyname", "The ec2 keypair to use").Default("default").String()
		poll      = kingpin.Command("poll", "Poll a cloudformation stack")
		pollStack = poll.Arg("stack", "The name of the stack to poll").Required().String()
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

	cwd, err := os.Getwd()
	if err != nil {
		ui.Fatal(err)
	}

	switch cmd {
	case up.FullCommand():
		cmds.UpCommand(ui, cmds.UpCommandInput{
			StackName: *upStack,
			Templates: core.Templates{core.FS(false)},
			Dir:       cwd,
			Cluster: cmds.EcsClusterParameters{
				KeyName: *upKeyName,
			},
		})
	case poll.FullCommand():
		cmds.PollCommand(ui, cmds.PollCommandInput{
			StackName: *pollStack,
		})
	}
}
