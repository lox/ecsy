package api

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var DefaultServices Services

type Services struct {
	Cloudformation cfnInterface
	ECS            ecsInterface
	Logs           cloudwatchLogsInterface
}

func init() {
	session := session.New(nil)
	DefaultServices.Cloudformation = cloudformation.New(session)
	DefaultServices.ECS = ecs.New(session)
	DefaultServices.Logs = cloudwatchlogs.New(session)
}
