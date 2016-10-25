package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/lox/ecsy/api"
	"github.com/lox/ecsy/templates"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	stackDateFormat = "20060102-150405"
)

func ConfigureCreateCluster(app *kingpin.Application, svc api.Services) {
	var cluster, keyName, instanceType, dockerUsername, dockerPassword, dockerEmail, authorizedKeys string
	var datadogKey, logspoutTarget string
	var instanceCount int

	cmd := app.Command("create-cluster", "Create an ECS cluster")
	cmd.Flag("cluster", "The name of the ECS cluster to create").
		Required().
		StringVar(&cluster)

	cmd.Flag("keyname", "The EC2 keypair to use for instance").
		Default("default").
		StringVar(&keyName)

	cmd.Flag("type", "The EC2 instance type to use").
		Default("t2.micro").
		StringVar(&instanceType)

	cmd.Flag("count", "The number of instances to use").
		Default("3").
		IntVar(&instanceCount)

	cmd.Flag("docker-username", "The docker Username to use").
		StringVar(&dockerUsername)

	cmd.Flag("docker-password", "The docker Password to use").
		StringVar(&dockerPassword)

	cmd.Flag("docker-email", "The docker Email to use").
		StringVar(&dockerEmail)

	cmd.Flag("datadog-key", "The datadog api key").
		StringVar(&datadogKey)

	cmd.Flag("logspout-target", "The endpoint to push logspout output to").
		StringVar(&logspoutTarget)

	cmd.Flag("authorized-keys", "A URL to fetch a SSH authorized_keys file from.").
		StringVar(&authorizedKeys)

	cmd.Action(func(c *kingpin.ParseContext) error {
		fmt.Printf("Creating cluster %s", cluster)

		_, err := svc.ECS.CreateCluster(&ecs.CreateClusterInput{
			ClusterName: aws.String(cluster),
		})

		network, err := getOrCreateNetworkStack(cluster, svc)
		if err != nil {
			return err
		}

		timer := time.Now()
		stackName := cluster + "-ecs-" + time.Now().Format(stackDateFormat)
		fmt.Printf("Creating cloudformation stack %s", stackName)

		err = api.CreateStack(svc.Cloudformation, stackName, templates.EcsStack(), map[string]string{
			"Subnets":            network.Subnets,
			"SecurityGroup":      network.SecurityGroup,
			"KeyName":            keyName,
			"ECSCluster":         cluster,
			"InstanceType":       instanceType,
			"DesiredCapacity":    strconv.Itoa(instanceCount),
			"DockerHubUsername":  dockerUsername,
			"DockerHubPassword":  dockerPassword,
			"DockerHubEmail":     dockerEmail,
			"LogspoutTarget":     logspoutTarget,
			"DatadogApiKey":      datadogKey,
			"AuthorizedUsersUrl": authorizedKeys,
		})
		if err != nil {
			return err
		}

		err = api.PollUntilCreated(svc.Cloudformation, stackName, func(event *cloudformation.StackEvent) {
			fmt.Printf("%s\n", api.FormatStackEvent(event))
		})
		if err != nil {
			return err
		}

		fmt.Printf("Cluster %s created in %s\n\n", cluster, time.Now().Sub(timer).String())
		return nil
	})
}

func getOrCreateNetworkStack(clusterName string, svc api.Services) (api.NetworkOutputs, error) {
	outputs, err := api.FindNetworkStack(svc.Cloudformation, clusterName)
	if err == nil {
		return outputs, nil
	}

	timer := time.Now()
	fmt.Printf("Creating Network Stack for %s", clusterName)

	err = api.CreateStack(svc.Cloudformation, outputs.StackName, templates.NetworkStack(), map[string]string{})
	if err != nil {
		return api.NetworkOutputs{}, err
	}

	err = api.PollUntilCreated(svc.Cloudformation, outputs.StackName, func(event *cloudformation.StackEvent) {
		fmt.Printf("%s\n", api.FormatStackEvent(event))
	})
	if err != nil {
		return api.NetworkOutputs{}, err
	}

	outputs, err = api.FindNetworkStack(svc.Cloudformation, clusterName)
	if err != nil {
		return api.NetworkOutputs{}, err
	}

	fmt.Printf("%s created in %s\n\n", outputs.StackName, time.Now().Sub(timer).String())
	return outputs, nil
}
