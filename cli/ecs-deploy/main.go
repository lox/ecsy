package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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

func main() {
	var (
		debug       = kingpin.Flag("debug", "Show debugging output").Bool()
		projectName = kingpin.Flag("project-name", "The name of the project").Short('p').Default(currentDirName()).String()
		composeFile = kingpin.Flag("file", "The docker-compose file to use").Short('f').Default("docker-compose.yml").String()
		cluster     = kingpin.Flag("cluster", "The ECS cluster to use").Short('c').Default("default").String()
		imageTags   = kingpin.Arg("imagetags", "Tags in the form image=tag to apply to the task").String()
	)

	kingpin.Version(Version)
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate)
	kingpin.CommandLine.Help =
		`Deploy updated task definitions to ECS`

	kingpin.Parse()
	ui := cli.DefaultUi

	if *debug == false {
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

	// look for the service stack
	log.Printf("Looking for service stack")
	serviceStack, err := cfnClient.FindStackByOutputs(map[string]string{
		"StackType":  "ecs-former::ecs-service",
		"ECSCluster": *cluster,
		"TaskFamily": *projectName,
	})
	if err != nil {
		ui.Fatal(err)
	}

	imageMap, err := parseImageMap(*imageTags)
	if err != nil {
		ui.Fatal(err)
	}

	ui.Println("Updating task definition with ECS")
	taskDef, err := ecsClient.UpdateComposerTaskDefinition(*composeFile, *projectName, imageMap)
	if err != nil {
		ui.Fatal(err)
	}

	serviceOutputs := serviceStack.OutputMap()

	err = ecsClient.UpdateService(serviceOutputs["ECSCluster"], serviceOutputs["ECSService"], taskDef.Arn)
	if err != nil {
		ui.Fatal(err)
	}

	ui.Printf("Waiting for service to stabilize")
	if err = ecsClient.WaitUntilServicesStable(*cluster, serviceOutputs["ECSService"]); err != nil {
		ui.Fatal(err)
	}

	ui.Println("Service available at", serviceOutputs["ECSLoadBalancer"])
}

func parseImageMap(s string) (ecs.ContainerImageMap, error) {
	m := ecs.ContainerImageMap{}

	for _, token := range strings.Split(s, ",") {
		pieces := strings.SplitN(token, "=", 2)
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
