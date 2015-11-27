package ecs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	amazonecs "github.com/aws/aws-sdk-go/service/ecs"
)

type Client struct {
	*amazonecs.ECS
}

func NewClient() *Client {
	return &Client{amazonecs.New(session.New(nil))}
}

func (c *Client) CreateCluster(name string) (*Cluster, error) {
	_, err := amazonecs.New(session.New()).CreateCluster(&amazonecs.CreateClusterInput{
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
