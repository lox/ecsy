package cmds

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/99designs/ecs-former/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type UpCommandInput struct {
	Templates core.Templates
	Dir       string
	Cluster   EcsClusterParameters
	StackName string
}

func UpCommand(ui *Ui, input UpCommandInput) {
	cfg, err := core.LoadConfig()
	if os.IsNotExist(err) {
		ui.Debug.Println("Config doesn't exist, creating one")

		stackName := input.StackName

		if stackName == "" {
			stackName = stackNameFromDir(input.Dir)

			if err = createEcsCluster(stackName, input.Templates, input.Cluster); err != nil {
				ui.Fatal(err)
			}
		}

		events, errs := pollStack(cloudformation.New(session.New()), stackName)

		for event := range events {
			ui.Println(eventString(event))
		}

		if err := <-errs; err != nil {
			ui.Fatal(err)
		}
	} else if err != nil {
		ui.Fatal(err)
	}

	log.Printf("%#v", cfg)
}

func stackNameFromDir(dir string) string {
	return path.Base(dir) + time.Now().Format("-20060102-150405")
}

type EcsClusterParameters struct {
	KeyName string
}

func createEcsCluster(stackName string, tpls core.Templates, params EcsClusterParameters) error {
	svc := cloudformation.New(session.New())
	tpl := tpls.EcsCluster()

	log.Printf("Creating ecs stack %s", stackName)
	resp, err := svc.CreateStack(&cloudformation.CreateStackInput{
		StackName: aws.String(stackName),
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
		},
		DisableRollback: aws.Bool(true),
		Parameters: []*cloudformation.Parameter{{
			ParameterKey:   aws.String("KeyName"),
			ParameterValue: aws.String(params.KeyName),
		}},
		TemplateBody: aws.String(string(tpl)),
	})

	log.Printf("Stack identifier is %s", *resp.StackId)
	return err
}
