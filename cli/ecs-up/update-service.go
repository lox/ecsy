package main

// type CreateServiceCommandInput struct {
// 	ClusterName       string
// 	DockerComposeFile string
// }

// func CreateServiceCommand(ui *cli.Ui, input CreateServiceCommandInput) {
// 	if serviceStack == nil {
// 		stackName = input.ProjectName + "-ecs-service-" + time.Now().Format("20060102-150405")
// 		ui.Printf("Failed to find an existing service stack, creating %s", stackName)

// 		_, err = cfnSvc.CreateStack(&cloudformation.CreateStackInput{
// 			StackName: aws.String(stackName),
// 			Capabilities: []*string{
// 				aws.String("CAPABILITY_IAM"),
// 			},
// 			Parameters: []*cloudformation.Parameter{{
// 				ParameterKey:   aws.String("ECSCluster"),
// 				ParameterValue: aws.String(input.ClusterName),
// 			}, {
// 				ParameterKey:   aws.String("TaskFamily"),
// 				ParameterValue: resp.TaskDefinition.Family,
// 			}, {
// 				ParameterKey:   aws.String("TaskDefinition"),
// 				ParameterValue: resp.TaskDefinition.TaskDefinitionArn,
// 			}, {
// 				ParameterKey:   aws.String("ContainerName"),
// 				ParameterValue: aws.String(input.ContainerName),
// 			}, {
// 				ParameterKey:   aws.String("Subnets"),
// 				ParameterValue: aws.String(subnets),
// 			}, {
// 				ParameterKey:   aws.String("Vpc"),
// 				ParameterValue: aws.String(vpc),
// 			}, {
// 				ParameterKey:   aws.String("ECSSecurityGroup"),
// 				ParameterValue: aws.String(securityGroup),
// 			}},
// 			TemplateBody: aws.String(string(input.Templates.EcsService())),
// 		})
// 		if err != nil {
// 			ui.Fatal(err)
// 		}
// 	} else {
// 		stackName := *serviceStack.StackName
// 		ui.Printf("Found existing service stack, updating %s", stackName)

// 		_, err := cfnSvc.UpdateStack(&cloudformation.UpdateStackInput{
// 			StackName: aws.String(stackName),
// 			Capabilities: []*string{
// 				aws.String("CAPABILITY_IAM"),
// 			},
// 			Parameters: []*cloudformation.Parameter{{
// 				ParameterKey:   aws.String("ECSCluster"),
// 				ParameterValue: aws.String(input.ClusterName),
// 			}, {
// 				ParameterKey:   aws.String("TaskFamily"),
// 				ParameterValue: resp.TaskDefinition.Family,
// 			}, {
// 				ParameterKey:   aws.String("TaskDefinition"),
// 				ParameterValue: resp.TaskDefinition.TaskDefinitionArn,
// 			}, {
// 				ParameterKey:   aws.String("ContainerName"),
// 				ParameterValue: aws.String(input.ContainerName),
// 			}, {
// 				ParameterKey:   aws.String("Subnets"),
// 				ParameterValue: aws.String(subnets),
// 			}, {
// 				ParameterKey:   aws.String("Vpc"),
// 				ParameterValue: aws.String(vpc),
// 			}, {
// 				ParameterKey:   aws.String("ECSSecurityGroup"),
// 				ParameterValue: aws.String(securityGroup),
// 			}},
// 			UsePreviousTemplate: aws.Bool(true),
// 		})
// 		if err != nil {
// 			ui.Fatal(err)
// 		}
// 	}

// 	events, errs := cloudformation.pollStack(cfnSvc, stackName)
// 	for event := range events {
// 		ui.Println(eventString(event))
// 	}

// 	if err := <-errs; err != nil {
// 		ui.Fatal(err)
// 	}
// }
