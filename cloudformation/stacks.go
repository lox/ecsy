package cloudformation

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	amazoncfn "github.com/aws/aws-sdk-go/service/cloudformation"
)

type Stack struct {
	StackName string
	stack     *amazoncfn.Stack
}

type OutputMap map[string]string

func (s *Stack) OutputMap() OutputMap {
	outputs := OutputMap{}
	for _, output := range s.stack.Outputs {
		outputs[*output.OutputKey] = *output.OutputValue
	}
	return outputs
}

func (o OutputMap) RequireKeys(keys ...string) error {
	for _, key := range keys {
		if _, ok := o[key]; !ok {
			return fmt.Errorf("Stack output is missing key %s", key)
		}
	}
	return nil
}

func (o OutputMap) Contains(match OutputMap) bool {
	for k, v := range match {
		if val, ok := o[k]; !ok || val != v {
			return false
		}
	}
	return true
}

var ErrNoStacksFound = errors.New("No matching stacks found")

func (c *Client) FindStackByOutputs(match OutputMap) (*Stack, error) {
	stacks, err := c.allActiveStacks()
	if err != nil {
		return nil, err
	}
	for _, stack := range stacks {
		if stack.OutputMap().Contains(match) {
			return stack, nil
		}
	}
	return nil, ErrNoStacksFound
}

func (c *Client) allActiveStacks() (stacks []*Stack, err error) {
	err = c.DescribeStacksPages(nil, func(page *amazoncfn.DescribeStacksOutput, last bool) bool {
		for _, s := range page.Stacks {
			if *s.StackStatus == "CREATE_COMPLETE" || *s.StackStatus == "UPDATE_COMPLETE" {
				stacks = append(stacks, &Stack{stack: s, StackName: *s.StackName})
			}
		}
		return last
	})
	return
}

type StackParameters map[string]string

func (sp StackParameters) slice() []*amazoncfn.Parameter {
	slice := []*amazoncfn.Parameter{}

	for k, v := range sp {
		slice = append(slice, &amazoncfn.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		})
	}

	return slice
}

func (c *Client) CreateStack(name string, body string, params StackParameters) error {
	_, err := c.CloudFormation.CreateStack(&amazoncfn.CreateStackInput{
		StackName: aws.String(name),
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
		},
		DisableRollback: aws.Bool(false),
		Parameters:      params.slice(),
		TemplateBody:    aws.String(body),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) StackOutputs(name string) (OutputMap, error) {
	resp, err := c.CloudFormation.DescribeStacks(&amazoncfn.DescribeStacksInput{
		StackName: aws.String(name),
	})

	if err != nil {
		return nil, err
	}

	outputs := OutputMap{}
	for _, output := range resp.Stacks[0].Outputs {
		outputs[*output.OutputKey] = *output.OutputValue
	}

	return outputs, nil
}
