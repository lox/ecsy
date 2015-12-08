package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/99designs/ecs-cli/cli"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var (
	Version string
	cfnSvc  *cloudformation.CloudFormation
	ecsSvc  *ecs.ECS
)

func main() {
	var (
		debug       = kingpin.Flag("debug", "Show debugging output").Bool()
		projectName = kingpin.Flag("project-name", "The name of the project").Short('p').Default(currentDirName()).String()
		composeFile = kingpin.Flag("file", "The docker-compose file to use").Short('f').Default("docker-compose.yml").String()
		cluster     = kingpin.Flag("cluster", "The ECS cluster to use").Short('c').Default("default").String()
		healthcheck = kingpin.Flag("healthcheck", "Path to the healthcheck route").Default("/").String()
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

	imageMap, err := parseImageMap(*imageTags)
	if err != nil {
		ui.Fatal(err)
	}

	if _, err := os.Stat(*composeFile); os.IsNotExist(err) {
		ui.Fatalf("Unable to open compose file %s", *composeFile)
	}

	session := session.New(nil)
	cfnSvc = cloudformation.New(session)
	ecsSvc = ecs.New(session)

	DeployCommand(ui, DeployCommandInput{
		ClusterName:     *cluster,
		ProjectName:     *projectName,
		ComposeFile:     *composeFile,
		HealthCheckUrl:  *healthcheck,
		ContainerImages: imageMap,
	})
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
