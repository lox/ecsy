package cmds

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type ListCommandInput struct {
}

func ListCommand(ui *Ui, input ListCommandInput) {
	svc := ecs.New(session.New())

	resp, err := svc.ListClusters(nil)
	if err != nil {
		ui.Fatal(err)
	}

	clusters := []*string{}

	for _, arn := range resp.ClusterArns {
		clusters = append(clusters, aws.String(strings.Split(*arn, "/")[1]))
	}

	descrResp, err := svc.DescribeClusters(&ecs.DescribeClustersInput{
		Clusters: clusters,
	})

	if err != nil {
		ui.Fatal(err)
	}

	// Pretty-print the response data.
	fmt.Println(descrResp)

}

// 	svc := cloudformation.New(session.New())

// 	resp, err := svc.ListStacks(&cloudformation.ListStacksInput{
// 		StackStatusFilter: []*string{
// 			aws.String("CREATE_COMPLETE"),
// 		},
// 	})

// 	if err != nil {
// 		ui.Fatal(err)
// 	}

// 	for _, stack := range resp.StackSummaries {
// 		ui.Printf("%s - %#v", *stack.StackName, isEcsStack(svc, *stack.StackName))
// 	}
// }

// func isEcsStack(svc *cloudformation.CloudFormation, stackName string) bool {
// 	params := &cloudformation.ListStackResourcesInput{
// 		StackName: aws.String(stackName),
// 	}

// 	resp, err := svc.ListStackResources(params)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Pretty-print the response data.
// 	fmt.Println(resp)

// 	return false
// }
