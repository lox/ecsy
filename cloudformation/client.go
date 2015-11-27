package cloudformation

import (
	"github.com/aws/aws-sdk-go/aws/session"
	amazoncfn "github.com/aws/aws-sdk-go/service/cloudformation"
)

type Client struct {
	*amazoncfn.CloudFormation
}

func NewClient() *Client {
	return &Client{amazoncfn.New(session.New())}
}
