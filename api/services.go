package api

import (
	"log"

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
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatal(err)
	}

	DefaultServices.Cloudformation = cloudformation.New(sess)
	DefaultServices.ECS = ecs.New(sess)
	DefaultServices.Logs = cloudwatchlogs.New(sess)
}
