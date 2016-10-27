package api

import (
	logs "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type cloudwatchLogsInterface interface {
	DescribeLogStreamsPages(input *logs.DescribeLogStreamsInput, fn func(p *logs.DescribeLogStreamsOutput, lastPage bool) (shouldContinue bool)) error
	FilterLogEventsPages(input *logs.FilterLogEventsInput, fn func(p *logs.FilterLogEventsOutput, lastPage bool) (shouldContinue bool)) error
}
