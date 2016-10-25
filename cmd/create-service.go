package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/lox/ecsy/api"
	"github.com/lox/ecsy/compose"
	"github.com/lox/ecsy/templates"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureCreateService(app *kingpin.Application, svc api.Services) {
	var cluster, projectName, healthCheck, certificateID string
	var composeFiles []string

	cmd := app.Command("create-service", "Create an ECS service for your app")
	cmd.Flag("cluster", "The name of the ECS cluster to use").
		Required().
		StringVar(&cluster)

	cmd.Flag("project-name", "The name of the Compose project").
		Short('p').
		Default(currentDirName()).
		StringVar(&projectName)

	cmd.Flag("healthcheck", "Path to check for HTTP health check").
		Default("/").
		StringVar(&healthCheck)

	cmd.Flag("ssl-certificate-id", "The identifier of the SSL certificate to associate with the service").
		Default("").
		StringVar(&certificateID)

	cmd.Flag("file", "The paths to docker-compose files to convert to service definitions").
		Short('f').
		Default("docker-compose.yml").
		ExistingFilesVar(&composeFiles)

	cmd.Action(func(c *kingpin.ParseContext) error {
		log.Printf("Creating service %s on %s", projectName, cluster)

		stack, _ := api.FindServiceStack(svc.Cloudformation, cluster, projectName)
		if stack != nil {
			return fmt.Errorf("A service already exists for %q in cluster %q. Use `deploy` or `update-service`",
				projectName, cluster)
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

		log.Printf("Registering a task for %s", projectName)
		resp, err := svc.ECS.RegisterTaskDefinition(taskDefinitionInput)
		if err != nil {
			return err
		}
		log.Printf("Registered task definition %s:%d", *resp.TaskDefinition.Family, *resp.TaskDefinition.Revision)

		network, err := api.FindNetworkStack(svc.Cloudformation, cluster)
		if err != nil {
			return err
		}
		log.Printf("Found network stack %s", network.StackName)

		params := map[string]string{
			"ECSCluster":       cluster,
			"TaskFamily":       *resp.TaskDefinition.Family,
			"TaskDefinition":   *resp.TaskDefinition.TaskDefinitionArn,
			"Subnets":          network.Subnets,
			"Vpc":              network.Vpc,
			"ECSSecurityGroup": network.SecurityGroup,
			"SSLCertificateId": certificateID,
		}

		exposedPorts := api.ExposedPorts(resp.TaskDefinition)

		if len(exposedPorts) != 1 {
			return fmt.Errorf("Task definition without exactly 1 host mapped port are not yet supported")
		}

		// for now this is a single value
		for container, mappings := range exposedPorts {
			for _, mapping := range mappings {
				params["ContainerName"] = container
				params["ContainerPort"] = strconv.FormatInt(*mapping.ContainerPort, 10)
				params["HealthCheckUrl"] = healthCheck
				params["ELBPort"] = strconv.FormatInt(*mapping.HostPort, 10)
			}
		}

		timer := time.Now()
		stackName := cluster + "-ecs-service-" + time.Now().Format(stackDateFormat)

		log.Printf("Creating service cloudformation stack %s", stackName)

		err = api.CreateStack(svc.Cloudformation, stackName, templates.EcsService(), params)
		if err != nil {
			return err
		}

		err = api.PollUntilCreated(svc.Cloudformation, stackName, func(event *cloudformation.StackEvent) {
			log.Printf("%s\n", api.FormatStackEvent(event))
		})
		if err != nil {
			return err
		}

		stackOutputs, err := api.StackOutputs(svc.Cloudformation, stackName)
		if err != nil {
			return err
		}

		var printer = func(e *ecs.ServiceEvent) {
			log.Println(*e.Message)
		}

		log.Printf("Waiting for service to reach a steady state.")
		err = api.PollUntilTaskDeployed(svc.ECS, cluster, stackOutputs["ECSService"], *resp.TaskDefinition.TaskDefinitionArn, printer)
		if err != nil {
			return err
		}

		// ui.Printf("Waiting for service to stabilize")
		// if err = apient.WaitUntilServicesStable(input.ClusterName, serviceOutputs["ECSService"]); err != nil {
		// 	ui.Fatal(err)
		// }

		log.Printf("Service created in %s", time.Now().Sub(timer).String())
		log.Printf("Service available at", stackOutputs["ECSLoadBalancer"])
		return nil
	})
}

func currentDirName() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Base(cwd)
}
