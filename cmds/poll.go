package cmds

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/mgutz/ansi"
)

type PollCommandInput struct {
	StackName string
}

func PollCommand(ui *Ui, input PollCommandInput) {
	events, errs := pollStack(cfn.New(session.New()), input.StackName)

	for event := range events {
		ui.Println(eventString(event))
	}

	if err := <-errs; err != nil {
		ui.Print(err)
		ui.Exit(1)
	}

}

func eventString(event *cfn.StackEvent) string {
	descr := ""
	if event.ResourceStatusReason != nil {
		descr = fmt.Sprintf("=> %q", *event.ResourceStatusReason)
	}
	return fmt.Sprintf("%-30s %s [%s] %s",
		colorizeStatus(*event.ResourceStatus),
		*event.LogicalResourceId,
		*event.ResourceType,
		descr,
	)
}

func colorizeStatus(s string) string {
	if strings.HasSuffix(s, "COMPLETE") {
		return ansi.Color(s, "green")
	} else if strings.HasSuffix(s, "FAILED") {
		return ansi.Color(s, "red")
	}
	return ansi.Color(s, "white")
}

func pollStack(svc *cfn.CloudFormation, stackName string) (chan *cfn.StackEvent, chan error) {
	poller := stackPoller{
		stackName: stackName,
		ch:        make(chan *cfn.StackEvent),
	}

	errs := make(chan error)
	go func() {
		params := &cfn.DescribeStackEventsInput{
			StackName: aws.String(stackName),
		}

		for !poller.done {
			events := events{}
			svc.DescribeStackEventsPages(params, func(page *cfn.DescribeStackEventsOutput, lastPage bool) bool {
				for _, ev := range page.StackEvents {
					if !events.Contains(ev) {
						events = append(events, ev)
					} else {
						return true
					}
				}
				return lastPage
			})
			poller.Push(events...)
			if poller.Error() != nil {
				break
			}
		}
		poller.Close()
		if err := poller.Error(); err != nil {
			errs <- err
		}
		close(errs)
	}()

	return poller.ch, errs
}

type stackPoller struct {
	ch        chan *cfn.StackEvent
	stackName string
	done      bool
	lastSeen  time.Time
	err       error
}

func (s *stackPoller) Push(events ...*cfn.StackEvent) {
	sort.Sort(eventTimeSorter(events))
	for _, ev := range events {
		if !ev.Timestamp.After(s.lastSeen) {
			continue
		}
		s.lastSeen = *ev.Timestamp
		s.ch <- ev
		if *ev.LogicalResourceId == s.stackName {
			switch *ev.ResourceStatus {
			case cfn.ResourceStatusCreateInProgress:
			case cfn.ResourceStatusUpdateFailed:
				fallthrough
			case cfn.ResourceStatusCreateFailed:
				s.err = errors.New(*ev.ResourceStatusReason)
				fallthrough
			default:
				s.done = true
				return
			}
		}
	}
}

func (s *stackPoller) Error() error {
	return s.err
}

func (s *stackPoller) Close() {
	s.done = true
	close(s.ch)
}

type events []*cfn.StackEvent

func (e events) Contains(ev *cfn.StackEvent) bool {
	for _, ei := range e {
		if ei.EventId == ev.EventId {
			return true
		}
	}
	return false
}

type eventTimeSorter []*cfn.StackEvent

func (e eventTimeSorter) Len() int           { return len(e) }
func (e eventTimeSorter) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e eventTimeSorter) Less(i, j int) bool { return e[i].Timestamp.Before(*e[j].Timestamp) }
