package api

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"time"
)

const ECS_POLL_INTERVAL = 1 * time.Second

type ecsInterface interface {
	DescribeServices(*ecs.DescribeServicesInput) (*ecs.DescribeServicesOutput, error)
	CreateCluster(*ecs.CreateClusterInput) (*ecs.CreateClusterOutput, error)
}

func UpdateContainerImages(defs []*ecs.ContainerDefinition, images map[string]string) error {
	for name, image := range images {
		var matched bool
		for idx, containerDef := range defs {
			if *containerDef.Name == name {
				matched = true

				wanted := parseDockerImageName(image)
				current := parseDockerImageName(*containerDef.Image)

				log.Println(wanted, current)

				if wanted.Image != "" {
					defs[idx].Image = aws.String(wanted.String())
				} else {
					current.Tag = wanted.Tag
					defs[idx].Image = aws.String(current.String())
				}
			}
		}
		if !matched {
			return fmt.Errorf("No defined container named %q", name)
		}
	}

	return nil
}

type dockerImageName struct {
	Image string
	Tag   string
}

func parseDockerImageName(image string) *dockerImageName {
	parts := strings.SplitN(image, ":", 2)
	if len(parts) == 1 {
		parts = append(parts, "latest")
	}

	return &dockerImageName{Image: parts[0], Tag: parts[1]}
}

func (d *dockerImageName) String() string {
	return fmt.Sprintf("%s:%s", d.Image, d.Tag)
}

func getService(svc ecsInterface, cluster, service string) (*ecs.Service, error) {
	resp, err := svc.DescribeServices(&ecs.DescribeServicesInput{
		Services: []*string{aws.String(service)},
		Cluster:  aws.String(cluster),
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Failures) > 0 {
		return nil, errors.New(*resp.Failures[0].Reason)
	}

	if len(resp.Services) != 1 {
		return nil, errors.New("Multiple services with the same name.")
	}

	return resp.Services[0], nil
}

func PollUntilTaskDeployed(svc ecsInterface, cluster string, service string, task string, f func(e *ecs.ServiceEvent)) error {
	lastSeen := time.Now().Add(-1 * time.Minute)

	for {
		service, err := getService(svc, cluster, service)
		if err != nil {
			return err
		}

		for i := len(service.Events) - 1; i >= 0; i-- {
			event := service.Events[i]
			if event.CreatedAt.After(lastSeen) {
				f(event)
				lastSeen = *event.CreatedAt
			}
		}

		if len(service.Deployments) == 1 && *service.Deployments[0].TaskDefinition == task {
			return nil
		}

		time.Sleep(ECS_POLL_INTERVAL)
	}
}

func ExposedPorts(taskDef *ecs.TaskDefinition) map[string][]*ecs.PortMapping {
	mappings := map[string][]*ecs.PortMapping{}

	for _, container := range taskDef.ContainerDefinitions {
		for _, mapping := range container.PortMappings {
			if mapping.HostPort != nil {
				if values, exists := mappings[*container.Name]; exists {
					mappings[*container.Name] = append(values, mapping)
				} else {
					mappings[*container.Name] = []*ecs.PortMapping{mapping}
				}
			}
		}
	}

	return mappings
}
