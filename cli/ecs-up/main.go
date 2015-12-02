package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/99designs/ecs-cli/cli"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ecs"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	Version string
	cfnSvc  *cloudformation.CloudFormation
	ecsSvc  *ecs.ECS
)

func main() {
	var (
		debug = kingpin.Flag("debug", "Show debugging output").Bool()

		// create-cluster command
		create                  = kingpin.Command("create-cluster", "Create an ECS cluster")
		createCluster           = create.Flag("cluster", "The name of the ECS cluster to create").Required().String()
		createKeyName           = create.Flag("keyname", "The EC2 keypair to use for instance").Default("default").String()
		createInstanceType      = create.Flag("type", "The EC2 instance type to use").Default("t2.micro").String()
		createInstanceCount     = create.Flag("count", "The number of instances to use").Default("3").Int()
		createDockerUsername    = create.Flag("docker-username", "The docker Username to use").String()
		createDockerPassword    = create.Flag("docker-password", "The docker Password to use").String()
		createDockerEmail       = create.Flag("docker-email", "The docker Email to use").String()
		createDatadogKey        = create.Flag("datadog-key", "The datadog api key").String()
		createLogspoutTarget    = create.Flag("logspout-target", "The endpoint to push logspout output to").String()
		createAuthorizedKeysUrl = create.Flag("authorized-keys", "A URL to fetch a SSH authorized_keys file from.").String()

		// create-service command
		createSvc            = kingpin.Command("create-service", "Create an ECS service for your app")
		createSvcCluster     = createSvc.Flag("cluster", "The ECS cluster to use").Required().String()
		createSvcProjectName = createSvc.Flag("project-name", "The name of the compose project").Short('p').Default(currentDirName()).String()
		createSvcFile        = createSvc.Flag("file", "The path to the docker-compose.yml file").Short('f').Default("docker-compose.yml").String()

		update = kingpin.Command("update-cluster", "Updates an existing ECS cluster")
		// updateName    = update.Arg("name", "The name of the ECS stack to create").Required().String()
		// updateKeyName = update.Flag("keyname", "The EC2 keypair to use for instance").Default("default").String()
		poll      = kingpin.Command("poll", "Poll a cloudformation stack")
		pollStack = poll.Arg("stack", "The name of the cloudformation stack to poll").String()
	)

	session := session.New(nil)
	cfnSvc = cloudformation.New(session)
	ecsSvc = ecs.New(session)

	kingpin.Version(Version)
	kingpin.CommandLine.Help =
		`Create ECS Cluster, Tasks and Services`

	ui := cli.DefaultUi
	cmd := kingpin.Parse()

	if *debug == false {
		log.SetOutput(ioutil.Discard)
	}

	switch cmd {
	case create.FullCommand():
		CreateClusterCommand(ui, CreateClusterCommandInput{
			ClusterName: *createCluster,
			Parameters: EcsClusterParameters{
				KeyName:            *createKeyName,
				InstanceType:       *createInstanceType,
				InstanceCount:      *createInstanceCount,
				DockerHubUsername:  *createDockerUsername,
				DockerHubPassword:  *createDockerPassword,
				DockerHubEmail:     *createDockerEmail,
				DatadogApiKey:      *createDatadogKey,
				LogspoutTarget:     *createLogspoutTarget,
				AuthorizedUsersUrl: *createAuthorizedKeysUrl,
			},
		})

	case createSvc.FullCommand():
		CreateServiceCommand(ui, CreateServiceCommandInput{
			ClusterName: *createSvcCluster,
			ProjectName: *createSvcProjectName,
			ComposeFile: *createSvcFile,
		})

	case update.FullCommand():
		UpdateClusterCommand(ui, UpdateClusterCommandInput{
		// Name:      *updateName,
		// Templates: templates.Templates{core.FS(false)},
		// Parameters: EcsClusterParameters{
		// 	KeyName: *updateKeyName,
		// },
		})

	case poll.FullCommand():
		PollCommand(ui, PollCommandInput{
			StackName: *pollStack,
		})
	}
}

func currentDirName() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Base(cwd)
}
