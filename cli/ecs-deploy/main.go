package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/99designs/ecs-cli/cli"
	"github.com/99designs/ecs-cli/cloudformation"
	"github.com/99designs/ecs-cli/ecs"
)

var (
	cfnClient *cloudformation.Client
	ecsClient *ecs.Client
	Version   string
)

func defaultProjectName() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Base(cwd)
}

func main() {
	var (
		debug       = kingpin.Flag("debug", "Show debugging output").Bool()
		projectName = kingpin.Flag("project-name", "The name of the project").Short('p').Default(defaultProjectName()).String()
		composeFile = kingpin.Flag("file", "The docker-compose file to use").Short('f').Default("docker-compose.yml").String()
		cluster     = kingpin.Flag("cluster", "The ECS cluster to use").Short('c').Default("default").String()
	)

	kingpin.Version(Version)
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate)
	kingpin.CommandLine.Help =
		`Deploy updated task definitions to ECS`

	kingpin.Parse()
	ui := cli.DefaultUi

	if *debug {
		log.SetOutput(ioutil.Discard)
	}

	if cfnClient == nil {
		cfnClient = cloudformation.NewClient()
	}

	clusterStack, err := cfnClient.FindStackByOutputs(map[string]string{
		"StackType":  "ecs-former::ecs-stack",
		"ECSCluster": *cluster,
	})
	if err != nil {
		ui.Fatal(err)
	}

	log.Printf("Found cluster stack %s", *cluster)
	outputs := clusterStack.OutputMap()

	if err = outputs.RequireKeys("Subnets", "Vpc", "SecurityGroup"); err != nil {
		ui.Fatal(err)
	}

	if ecsClient == nil {
		ecsClient = ecs.NewClient()
	}

	ui.Println("Registering task definition with ECS")
	taskDef, err := ecsClient.RegisterComposerTaskDefinition(*composeFile, *projectName)
	if err != nil {
		ui.Fatal(err)
	}

	log.Printf("Registered task definition %s:%d (%s)", taskDef.Family, taskDef.Revision, taskDef.Arn)

	// look for the service stack
	serviceStack, err := cfnClient.FindStackByOutputs(map[string]string{
		"StackType":  "ecs-former::ecs-service",
		"ECSCluster": *cluster,
		"TaskFamily": taskDef.Family,
	})
	if err != nil {
		ui.Fatal(err)
	}

	log.Printf("%#v", serviceStack)
}
