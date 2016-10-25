package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/lox/ecsy/api"
	"github.com/lox/ecsy/compose"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureDeploy(app *kingpin.Application, svc api.Services) {
	var cluster, projectName, imageTags string
	var composeFiles []string

	cmd := app.Command("deploy", "Deploy updated task definitions to ECS")
	cmd.Flag("cluster", "The ECS cluster to deploy to").
		Required().
		StringVar(&cluster)

	cmd.Flag("project-name", "The name of the project").
		Short('p').
		Default(currentDirName()).
		StringVar(&projectName)

	cmd.Flag("file", "The docker-compose file to use").
		Short('f').
		Default("docker-compose.yml").
		ExistingFilesVar(&composeFiles)

	cmd.Arg("imagetags", "Tags in the form image=tag to apply to the task").
		StringVar(&imageTags)

	cmd.Action(func(c *kingpin.ParseContext) error {
		images, err := parseImageMap(imageTags)
		if err != nil {
			return err
		}

		log.Printf("Generating task definition from %#v", composeFiles)
		t := compose.Transformer{
			ComposeFiles: composeFiles,
			ProjectName:  projectName,
		}

		taskDefinitionInput, err := t.Transform()
		if err != nil {
			return err
		}

		log.Printf("Updating task definition for task %s", *taskDefinitionInput.Family)
		err = api.UpdateContainerImages(taskDefinitionInput.ContainerDefinitions, images)
		if err != nil {
			return err
		}

		resp, err := svc.ECS.RegisterTaskDefinition(taskDefinitionInput)
		if err != nil {
			return err
		}
		log.Printf("Registered task definition %s:%d", *resp.TaskDefinition.Family, *resp.TaskDefinition.Revision)

		serviceStack, err := api.FindServiceStack(svc.Cloudformation, cluster, projectName)
		if err != nil {
			return err
		}
		log.Printf("Found service stack %s", *serviceStack.StackName)

		outputs := api.StackOutputMap(serviceStack)
		timer := time.Now()

		log.Printf("Updating service %s with new task definition", *serviceStack.StackName)
		_, err = svc.ECS.UpdateService(&ecs.UpdateServiceInput{
			Service:        aws.String(outputs["ECSService"]),
			Cluster:        aws.String(outputs["ECSCluster"]),
			TaskDefinition: aws.String(*resp.TaskDefinition.TaskDefinitionArn),
		})
		if err != nil {
			return err
		}

		var printer = func(e *ecs.ServiceEvent) {
			log.Println(*e.Message)
		}

		log.Printf("Waiting for service to reach a steady state.")
		err = api.PollUntilTaskDeployed(svc.ECS, outputs["ECSCluster"], outputs["ECSService"], *resp.TaskDefinition.TaskDefinitionArn, printer)
		if err != nil {
			return err
		}

		// ui.Printf("Waiting for service to stabilize")
		// if err = apient.WaitUntilServicesStable(input.ClusterName, serviceOutputs["ECSService"]); err != nil {
		// 	ui.Fatal(err)
		// }

		log.Printf("Deployed in %s", time.Now().Sub(timer).String())
		return nil
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
