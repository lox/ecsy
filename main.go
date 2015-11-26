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
		debug             = kingpin.Flag("debug", "Show debugging output").Bool()
		list              = kingpin.Command("ls", "List ecs-former stacks")
		create            = kingpin.Command("create-cluster", "Create an ECS cluster")
		createName        = create.Arg("name", "The name of the ECS stack to create").Required().String()
		createKeyName     = create.Flag("keyname", "The EC2 keypair to use for instance").Default("default").String()
		update            = kingpin.Command("update-cluster", "Updates an existing ECS cluster")
		updateName        = update.Arg("name", "The name of the ECS stack to create").Required().String()
		updateKeyName     = update.Flag("keyname", "The EC2 keypair to use for instance").Default("default").String()
		poll              = kingpin.Command("poll", "Poll a CloudFormation stack")
		pollStack         = poll.Arg("stack", "The name of the stack to poll").Required().String()
		deploy            = kingpin.Command("deploy", "Deploy an app on ecs")
		deployComposeFile = deploy.Flag("file", "The docker-compose.yml file to use").Default("docker-compose.yml").String()
		deployProject     = deploy.Flag("project", "The name of the project").String()
		deployContainer   = deploy.Flag("container", "The name of the container to bind to a service").Required().String()
		deployCluster     = deploy.Flag("cluster", "The ecs cluster to deploy to").Required().String()
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
	case update.FullCommand():
		cmds.UpdateClusterCommand(ui, cmds.UpdateClusterCommandInput{
			Name:      *updateName,
			Templates: core.Templates{core.FS(false)},
			Parameters: cmds.EcsClusterParameters{
				KeyName: *updateKeyName,
			},
		})
	case poll.FullCommand():
		cmds.PollCommand(ui, cmds.PollCommandInput{
			StackName: *pollStack,
		})
	case list.FullCommand():
		cmds.ListCommand(ui, cmds.ListCommandInput{})
	case deploy.FullCommand():
		cmds.DeployCommand(ui, cmds.DeployCommandInput{
			ComposeFile:   *deployComposeFile,
			ProjectName:   *deployProject,
			ClusterName:   *deployCluster,
			ContainerName: *deployContainer,
			Templates:     core.Templates{core.FS(false)},
		})

	}
}
