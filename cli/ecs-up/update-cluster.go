package main

import "github.com/99designs/ecs-cli/cli"

type UpdateClusterCommandInput struct {
	Name       string
	Parameters EcsClusterParameters
}

func UpdateClusterCommand(ui *cli.Ui, input UpdateClusterCommandInput) {

	// ui.Printf("Updating cloudformation stack %s", input.StackName)
	// svc := cloudformation.New(session.New())
	// _, err = svc.CreateStack(&cloudformation.CreateStackInput{
	// 	StackName: aws.String(input.StackName),
	// 	Capabilities: []*string{
	// 		aws.String("CAPABILITY_IAM"),
	// 	},
	// 	DisableRollback: aws.Bool(true),
	// 	Parameters: []*cloudformation.Parameter{{
	// 		ParameterKey:   aws.String("KeyName"),
	// 		ParameterValue: aws.String(input.Parameters.KeyName),
	// 	}, {
	// 		ParameterKey:   aws.String("ECSCluster"),
	// 		ParameterValue: aws.String(input.Name),
	// 	}},
	// 	TemplateBody: aws.String(string(input.Templates.EcsStack())),
	// })
	// if err != nil {
	// 	ui.Fatal(err)
	// }

	// events, errs := pollStack(svc, input.StackName)
	// for event := range events {
	// 	ui.Println(eventString(event))
	// }

	// if err := <-errs; err != nil {
	// 	ui.Fatal(err)
	// }
}
