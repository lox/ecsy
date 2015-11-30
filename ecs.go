package ecscli

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

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

func PollDeployment(svc ecsInterface, service *ecs.Service, cluster string, f func(e *ecs.ServiceEvent)) error {
	if len(service.Deployments) == 0 {
		return errors.New("Service has no deployments")
	}

	deployment := service.Deployments[0]
	lastSeen := *deployment.CreatedAt

	for {
		for i := len(service.Events) - 1; i >= 0; i-- {
			if service.Events[i].CreatedAt.After(lastSeen) {
				f(service.Events[i])
				lastSeen = *service.Events[i].CreatedAt
			}
		}

		if *deployment.DesiredCount == *deployment.RunningCount &&
			strings.HasSuffix(*service.Events[0].Message, "reached a steady state.") {
			break
		}

		resp, err := svc.DescribeServices(&ecs.DescribeServicesInput{
			Services: []*string{service.ServiceName},
			Cluster:  aws.String(cluster),
		})
		if err != nil {
			return err
		}

		if len(resp.Failures) > 0 {
			return errors.New(*resp.Failures[0].Reason)
		}

		service = resp.Services[0]
		deployment = deploymentById(*deployment.Id, resp.Services[0].Deployments)

		if deployment == nil {
			return errors.New("Failed to find deployment " + *deployment.Id)
		}
	}

	return nil
}

func deploymentById(id string, deployments []*ecs.Deployment) *ecs.Deployment {
	for _, d := range deployments {
		if id == *d.Id {
			return d
		}
	}
	return nil
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
