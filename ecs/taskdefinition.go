package ecs

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	amazonecs "github.com/aws/aws-sdk-go/service/ecs"
)

type PortMapping struct {
	ContainerName           string
	ContainerPort, HostPort int64
	Protocol                string
}

type TaskDefinition struct {
	Family, Arn, Status string
	Revision            int64
	PortMappings        []PortMapping
}

func (c *Client) RegisterComposerTaskDefinition(file string, projectName string) (*TaskDefinition, error) {
	input, err := TransformComposeFile(file, projectName)
	if err != nil {
		return nil, err
	}

	log.Printf("Task Definition: %#v", input)

	return c.registerTaskDefinition(input)

}

type ContainerImageMap map[string]string

func (m ContainerImageMap) apply(defs []*amazonecs.ContainerDefinition) error {
	for name, tag := range m {
		var matched bool
		for idx, containerDef := range defs {
			if *containerDef.Name == name {
				matched = true
				imageParts := strings.Split(*containerDef.Image, ":")
				defs[idx].Image = aws.String(imageParts[0] + ":" + tag)
			}
		}
		if !matched {
			return fmt.Errorf("Failed to find a container named %s", name)
		}
	}

	return nil
}

func (c *Client) UpdateComposerTaskDefinition(file string, projectName string, images ContainerImageMap) (*TaskDefinition, error) {
	input, err := TransformComposeFile(file, projectName)
	if err != nil {
		return nil, err
	}

	if err = images.apply(input.ContainerDefinitions); err != nil {
		return nil, err
	}

	log.Printf("Task Definition: %#v", input)
	return c.registerTaskDefinition(input)
}

func (c *Client) registerTaskDefinition(input *amazonecs.RegisterTaskDefinitionInput) (*TaskDefinition, error) {
	resp, err := c.client.RegisterTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	mappings := []PortMapping{}

	for _, container := range resp.TaskDefinition.ContainerDefinitions {
		for _, mapping := range container.PortMappings {
			if mapping.HostPort != nil {
				mappings = append(mappings, PortMapping{
					ContainerName: *container.Name,
					ContainerPort: *mapping.ContainerPort,
					HostPort:      *mapping.HostPort,
					Protocol:      *mapping.Protocol,
				})
			}
		}
	}

	return &TaskDefinition{
		Family:       *resp.TaskDefinition.Family,
		Revision:     *resp.TaskDefinition.Revision,
		Arn:          *resp.TaskDefinition.TaskDefinitionArn,
		Status:       *resp.TaskDefinition.Status,
		PortMappings: mappings,
	}, nil
}
