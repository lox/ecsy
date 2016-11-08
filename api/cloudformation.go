package api

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/fatih/color"
)

type cfnInterface interface {
	DescribeStacksPages(*cloudformation.DescribeStacksInput, func(*cloudformation.DescribeStacksOutput, bool) bool) error
	DescribeStackEventsPages(*cloudformation.DescribeStackEventsInput, func(*cloudformation.DescribeStackEventsOutput, bool) bool) error
	DescribeStacks(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error)
	CreateStack(*cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error)
	DeleteStack(*cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error)
}

type stackOutputMap map[string]string

func StackOutputMap(stack *cloudformation.Stack) stackOutputMap {
	outputs := stackOutputMap{}
	for _, output := range stack.Outputs {
		outputs[*output.OutputKey] = *output.OutputValue
	}
	return outputs
}

func (o stackOutputMap) RequireKeys(keys ...string) error {
	for _, key := range keys {
		if _, ok := o[key]; !ok {
			return fmt.Errorf("Stack output is missing key %s", key)
		}
	}
	return nil
}

func (o stackOutputMap) Contains(match map[string]string) bool {
	for k, v := range match {
		if val, ok := o[k]; !ok || val != v {
			return false
		}
	}
	return true
}

func GetStackOutputByKey(stack *cloudformation.Stack, key string) (string, bool) {
	for _, output := range stack.Outputs {
		if *output.OutputKey == key {
			return *output.OutputValue, true
		}
	}
	return "", false
}

var ErrNoStacksFound = errors.New("No matching stacks found")

func FindStacksByOutputs(svc cfnInterface, match map[string]string) ([]*cloudformation.Stack, error) {
	stacks, err := findAllActiveStacks(svc)
	if err != nil {
		return nil, err
	}
	filteredStacks := make([]*cloudformation.Stack, 0)

	for _, stack := range stacks {
		if StackOutputMap(stack).Contains(match) {
			filteredStacks = append(filteredStacks, stack)
		}
	}
	return filteredStacks, nil
}

func findAllActiveStacks(svc cfnInterface) (stacks []*cloudformation.Stack, err error) {
	err = svc.DescribeStacksPages(nil, func(page *cloudformation.DescribeStacksOutput, last bool) bool {
		for _, s := range page.Stacks {
			if *s.StackStatus != "DELETE_COMPLETE" {
				stacks = append(stacks, s)
			}
		}
		return last
	})
	return
}

func FindStacksByName(svc cfnInterface, stackName string) (stacks []*cloudformation.Stack, err error) {
	filter := &cloudformation.DescribeStacksInput{
		StackName: &stackName,
	}

	err = svc.DescribeStacksPages(filter, func(page *cloudformation.DescribeStacksOutput, last bool) bool {
		for _, s := range page.Stacks {
			if *s.StackStatus != "DELETE_COMPLETE" {
				stacks = append(stacks, s)
			}
		}
		return last
	})
	return
}

type CreateStackContext struct {
	Params          map[string]string
	DisableRollback bool
}

func CreateStack(svc cfnInterface, name string, body string, ctx CreateStackContext) error {
	paramsSlice := []*cloudformation.Parameter{}
	for k, v := range ctx.Params {
		paramsSlice = append(paramsSlice, &cloudformation.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		})
	}

	_, err := svc.CreateStack(&cloudformation.CreateStackInput{
		StackName: aws.String(name),
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
		},
		DisableRollback: aws.Bool(ctx.DisableRollback),
		Parameters:      paramsSlice,
		TemplateBody:    aws.String(body),
	})
	if err != nil {
		return err
	}
	return nil
}

func DeleteStack(svc cfnInterface, name string) error {
	_, err := svc.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: &name,
	})

	return err
}

func StackOutputs(svc cfnInterface, name string) (stackOutputMap, error) {
	resp, err := svc.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: aws.String(name),
	})

	if err != nil {
		return nil, err
	}

	outputs := stackOutputMap{}
	for _, output := range resp.Stacks[0].Outputs {
		outputs[*output.OutputKey] = *output.OutputValue
	}

	return outputs, nil
}

func PollUntilCreated(svc cfnInterface, stackName string, f func(e *cloudformation.StackEvent)) error {
	return PollStackEventsUntil(svc, stackName, isCreateUpdateComplete, f)
}

func PollUntilDeleted(svc cfnInterface, stackName string, f func(e *cloudformation.StackEvent)) error {
	return PollStackEventsUntil(svc, stackName, isDeleteComplete, f)
}

func PollStackEventsUntil(svc cfnInterface, stackName string, terminalCondition EventChecker, f func(e *cloudformation.StackEvent)) error {
	lastSeen := time.Time{}

	for {
		events, err := allStackEvents(svc, stackName, lastSeen)
		if err != nil {
			return err
		}

		for i := len(events) - 1; i >= 0; i-- {
			if events[i].Timestamp.After(lastSeen) {
				f(events[i])
				lastSeen = *events[i].Timestamp
			}
		}

		if len(events) > 0 {
			t, err := terminalCondition(stackName, events[0])
			if err != nil {
				return err
			}
			if t {
				break
			}
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

func allStackEvents(svc cfnInterface, stackName string, after time.Time) (events []*cloudformation.StackEvent, err error) {
	params := &cloudformation.DescribeStackEventsInput{
		StackName: aws.String(stackName),
	}

	err = svc.DescribeStackEventsPages(params, func(page *cloudformation.DescribeStackEventsOutput, last bool) bool {
		for _, event := range page.StackEvents {
			if !event.Timestamp.After(after) {
				return true
			}
			events = append(events, event)
		}
		return last
	})

	return
}

type EventChecker func(stackName string, ev *cloudformation.StackEvent) (bool, error)

func isCreateUpdateComplete(stackName string, ev *cloudformation.StackEvent) (bool, error) {
	if *ev.LogicalResourceId == stackName {
		switch *ev.ResourceStatus {
		case cloudformation.ResourceStatusUpdateComplete,
			cloudformation.ResourceStatusCreateComplete,
			cloudformation.ResourceStatusUpdateFailed,
			cloudformation.ResourceStatusCreateFailed,
			cloudformation.StackStatusRollbackComplete,
			cloudformation.StackStatusRollbackFailed,
			cloudformation.StackStatusUpdateRollbackComplete,
			cloudformation.StackStatusUpdateRollbackFailed:
			var err error
			if ev.ResourceStatusReason != nil {
				err = errors.New(*ev.ResourceStatusReason)
			}
			return true, err
		}
	}
	return false, nil
}

func isDeleteComplete(stackName string, ev *cloudformation.StackEvent) (bool, error) {
	if *ev.LogicalResourceId == stackName {
		switch *ev.ResourceStatus {
		case cloudformation.ResourceStatusDeleteComplete:
			return true, nil
		case cloudformation.ResourceStatusDeleteFailed:
			return true, errors.New(*ev.ResourceStatusReason)
		}
	}
	return false, nil
}

func formatResourceStatus(s string) string {
	switch {
	case strings.HasSuffix(s, "COMPLETE") && !strings.HasPrefix(s, "DELETE"):
		return color.GreenString(s)
	case strings.Contains(s, "FAILED") || strings.Contains(s, "ROLLBACK"):
		return color.RedString(s)
	case strings.HasSuffix(s, "IN_PROGRESS"):
		return color.YellowString(s)
	}
	return s
}

func FormatStackEvent(event *cloudformation.StackEvent) string {
	descr := ""
	if event.ResourceStatusReason != nil {
		descr = fmt.Sprintf("=> %q", *event.ResourceStatusReason)
	}
	return fmt.Sprintf("%s -> %s [%s] %s",
		formatResourceStatus(*event.ResourceStatus),
		*event.LogicalResourceId,
		*event.ResourceType,
		descr,
	)
}
