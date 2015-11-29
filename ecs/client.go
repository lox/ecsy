package ecs

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	amazonecs "github.com/aws/aws-sdk-go/service/ecs"
)

type Client struct {
	client *amazonecs.ECS
}

func NewClient() *Client {
	return &Client{amazonecs.New(session.New(nil))}
}

func (c *Client) CreateCluster(name string) (*Cluster, error) {
	_, err := c.client.CreateCluster(&amazonecs.CreateClusterInput{
		ClusterName: aws.String(name),
	})
	if err != nil {
		return nil, err
	}
	return &Cluster{Name: name}, nil
}

type Cluster struct {
	Name string
}

func (c *Client) WaitUntilServicesStable(clusterName string, serviceArn ...string) error {
	input := amazonecs.DescribeServicesInput{
		Services: []*string{},
		Cluster:  aws.String(clusterName),
	}

	for _, service := range serviceArn {
		input.Services = append(input.Services, aws.String(service))
	}

	return c.client.WaitUntilServicesStable(&input)
}

func (c *Client) UpdateService(cluster, service, taskArn string) error {
	resp, err := c.client.UpdateService(&amazonecs.UpdateServiceInput{
		Service:        aws.String(service),
		Cluster:        aws.String(cluster),
		TaskDefinition: aws.String(taskArn),
	})
	if err != nil {
		return err
	}

	log.Printf("%#v", resp)
	return nil
}
