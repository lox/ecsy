package cloudformation

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	amazoncfn "github.com/aws/aws-sdk-go/service/cloudformation"
)

type StackEvent struct {
	ResourceStatus       string
	ResourceType         string
	LogicalResourceId    string
	Timestamp            time.Time
	ResourceStatusReason string
}

func (event *StackEvent) String() string {
	descr := ""
	if event.ResourceStatusReason != "" {
		descr = fmt.Sprintf("=> %q", event.ResourceStatusReason)
	}
	return fmt.Sprintf("%s -> %s [%s] %s",
		event.ResourceStatus,
		event.LogicalResourceId,
		event.ResourceType,
		descr,
	)
}

type stackPoller struct {
	stackName string
	client    *Client
	err       error
}

func (p *stackPoller) Events(after time.Time) <-chan StackEvent {
	ch := make(chan StackEvent)

	go func() {
		defer close(ch)

		log.Printf("Polling stack events after %s", after.String())
		for {
			events, err := p.client.stackEvents(p.stackName, after)
			if err != nil {
				p.err = err
				break
			}

			// iterate in reverse order, so we emit in linear order
			for i := len(events) - 1; i >= 0; i-- {
				ch <- events[i]
				after = events[i].Timestamp
			}

			// check if the final event is a stack complete message
			if len(events) > 0 {
				done, err := p.isStackDone(events[0])
				if done {
					p.err = err
					break
				}
			}
		}
	}()
	return ch
}

func (p *stackPoller) isStackDone(ev StackEvent) (bool, error) {
	if ev.LogicalResourceId == p.stackName {
		switch ev.ResourceStatus {
		// successful
		case amazoncfn.ResourceStatusUpdateComplete:
			fallthrough
		case amazoncfn.ResourceStatusCreateComplete:
			return true, nil
		// error
		case amazoncfn.ResourceStatusUpdateFailed:
			fallthrough
		case amazoncfn.ResourceStatusCreateFailed:
			return true, errors.New(ev.ResourceStatusReason)
		}
	}
	return false, nil
}

func (p *stackPoller) Error() error {
	return p.err
}

func (c *Client) PollStackEvents(name string) *stackPoller {
	return &stackPoller{client: c, stackName: name}
}

func (c *Client) stackEvents(stackName string, after time.Time) (events []StackEvent, err error) {
	params := &amazoncfn.DescribeStackEventsInput{
		StackName: aws.String(stackName),
	}

	err = c.DescribeStackEventsPages(params, func(page *amazoncfn.DescribeStackEventsOutput, last bool) bool {
		for _, ev := range page.StackEvents {
			if !ev.Timestamp.After(after) {
				return true
			}
			event := StackEvent{
				ResourceStatus:    *ev.ResourceStatus,
				LogicalResourceId: *ev.LogicalResourceId,
				Timestamp:         *ev.Timestamp,
				ResourceType:      *ev.ResourceType,
			}
			if ev.ResourceStatusReason != nil {
				event.ResourceStatusReason = *ev.ResourceStatusReason
			}
			events = append(events, event)
		}
		return last
	})
	if err != nil {
		log.Println(err)
	}

	return
}
