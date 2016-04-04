package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ecs"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
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

		deleteCluster     = kingpin.Command("delete-cluster", "Remove a cluster and all running services on it")
		deleteClusterName = deleteCluster.Flag("cluster", "The name of the ECS cluster to create").Required().String()

		// create-service command
		createSvc            = kingpin.Command("create-service", "Create an ECS service for your app")
		createSvcCluster     = createSvc.Flag("cluster", "The ECS cluster to use").Required().String()
		createSvcProjectName = createSvc.Flag("project-name", "The name of the compose project").Short('p').Default(currentDirName()).String()
		createSvcFile        = createSvc.Flag("file", "The path to the docker-compose.yml file").Short('f').Default("docker-compose.yml").String()
		createSvcHealthcheck = createSvc.Flag("healthcheck", "Path to the healthcheck route").Default("/").String()

		// poll command
		poll      = kingpin.Command("poll", "Poll a cloudformation stack")
		pollStack = poll.Arg("stack", "The name of the cloudformation stack to poll").String()

		// deploy command
		deploy            = kingpin.Command("deploy", "Deploy updated task definitions to ECS")
		deployProjectName = deploy.Flag("project-name", "The name of the project").Short('p').Default(currentDirName()).String()
		deployComposeFile = deploy.Flag("file", "The docker-compose file to use").Short('f').Default("docker-compose.yml").String()
		deployCluster     = deploy.Flag("cluster", "The ECS cluster to use").Short('c').Default("default").String()
		deployHealthcheck = deploy.Flag("healthcheck", "Path to the healthcheck route").Default("/").String()
		deployImageTags   = deploy.Arg("imagetags", "Tags in the form image=tag to apply to the task").String()
	)

	session := session.New(nil)
	cfnSvc = cloudformation.New(session)
	ecsSvc = ecs.New(session)

	kingpin.Version(Version)
	kingpin.CommandLine.Help =
		`An unofficial set of commands for bootstrapping and working with ECS`

	ui := DefaultUi
	cmd := kingpin.Parse()

	if *debug == false {
		log.SetOutput(ioutil.Discard)
	}

	deployImageMap, err := parseImageMap(*deployImageTags)
	if err != nil {
		ui.Fatal(err)
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

	case deleteCluster.FullCommand():
		DeleteClusterCommand(ui, *deleteClusterName)

	case createSvc.FullCommand():
		CreateServiceCommand(ui, CreateServiceCommandInput{
			ClusterName:    *createSvcCluster,
			ProjectName:    *createSvcProjectName,
			ComposeFile:    *createSvcFile,
			HealthCheckUrl: *createSvcHealthcheck,
		})

	case poll.FullCommand():
		PollCommand(ui, PollCommandInput{
			StackName: *pollStack,
		})

	case deploy.FullCommand():
		DeployCommand(ui, DeployCommandInput{
			ClusterName:     *deployCluster,
			ProjectName:     *deployProjectName,
			ComposeFile:     *deployComposeFile,
			HealthCheckUrl:  *deployHealthcheck,
			ContainerImages: deployImageMap,
		})
	}
}

func parseImageMap(s string) (map[string]string, error) {
	m := map[string]string{}

	if s == "" {
		return m, nil
	}

	for _, token := range strings.Split(s, ",") {
		pieces := strings.SplitN(token, "=", 2)
		if len(pieces) != 2 {
			return nil, fmt.Errorf("Failed to parse image tag %q", s)
		}
		m[pieces[0]] = pieces[1]
	}

	log.Printf("Parsed %s to %#v", s, m)
	return m, nil

}

func currentDirName() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Base(cwd)
}
