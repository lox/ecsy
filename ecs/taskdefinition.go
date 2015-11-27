package ecs

import "log"

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

	resp, err := c.RegisterTaskDefinition(input)
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
