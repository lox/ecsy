package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/lox/ecsy/api"
	"github.com/lox/ecsy/taskdefinitions"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureDeploy(app *kingpin.Application, svc api.Services) {
	var cluster, serviceName, taskName, imageTags string
	var taskDefinitionFile string

	cmd := app.Command("deploy", "Deploy updated task definitions to ECS")
	cmd.Flag("cluster", "The ECS cluster to deploy to").
		Required().
		StringVar(&cluster)

	cmd.Flag("name", "The name of the task").
		Required().
		StringVar(&taskName)

	cmd.Flag("service", "The name of the ECS Service to deploy to").
		Required().
		StringVar(&serviceName)

	cmd.Flag("file", "The paths to task definitions in yaml or json").
		Short('f').
		ExistingFileVar(&taskDefinitionFile)

	cmd.Arg("imagetags", "Tags in the form image=tag to apply to the task").
		StringVar(&imageTags)

	cmd.Action(func(c *kingpin.ParseContext) error {
		images, err := parseImageMap(imageTags)
		if err != nil {
			return err
		}

		taskDefinitionInput, err := taskdefinitions.ParseFile(taskDefinitionFile, os.Environ())
		if err != nil {
			return err
		}

		clusterStack, err := api.FindClusterStack(svc.Cloudformation, cluster)
		if err != nil {
			return err
		} else if clusterStack == nil {
			return fmt.Errorf("No cluster exists for %q. Use `create-cluster`",
				cluster)
		}

		if logGroup, exists := api.GetStackOutputByKey(clusterStack, "LogGroupName"); exists {
			log.Printf("Setting tasks to use log group %s", logGroup)

			for _, def := range taskDefinitionInput.ContainerDefinitions {
				if def.LogConfiguration == nil {
					def.LogConfiguration = &ecs.LogConfiguration{
						LogDriver: aws.String("awslogs"),
						Options: map[string]*string{
							"awslogs-group":         aws.String(logGroup),
							"awslogs-region":        aws.String(os.Getenv("AWS_REGION")),
							"awslogs-stream-prefix": aws.String(serviceName),
						},
					}
				}
			}
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

		serviceStack, err := api.FindServiceStack(svc.Cloudformation, cluster, serviceName)
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
