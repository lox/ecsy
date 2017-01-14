package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	logs "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/lox/ecsy/api"
	"github.com/lox/ecsy/compose"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureRunTask(app *kingpin.Application, svc api.Services) {
	var cluster, projectName, service string
	var composeFiles []string
	var commands []string

	cmd := app.Command("run-task", "Run a once-off task")
	cmd.Alias("run")

	cmd.Flag("cluster", "The name of the ECS cluster to use").
		Required().
		StringVar(&cluster)

	cmd.Flag("project-name", "The name of the Compose project").
		Short('p').
		Default(currentDirName()).
		StringVar(&projectName)

	cmd.Flag("file", "The paths to docker-compose files to convert to task definitions").
		Short('f').
		Default("docker-compose.yml").
		ExistingFilesVar(&composeFiles)

	cmd.Flag("service", "The name of compose service to run").
		Short('s').
		StringVar(&service)

	cmd.Arg("commands", "Commands to override the default task command with").
		StringsVar(&commands)

	cmd.Action(func(c *kingpin.ParseContext) error {
		taskName := fmt.Sprintf("%s_%s_run", projectName, service)
		log.Printf("Creating task %s on %s", taskName, cluster)

		clusterStack, err := api.FindClusterStack(svc.Cloudformation, cluster)
		if err != nil {
			return err
		}
		if clusterStack == nil {
			return fmt.Errorf("No cluster exists for %q. Use `create-cluster`",
				cluster)
		}

		log.Printf("Generating task definition from %v", composeFiles)
		t := compose.Transformer{
			ComposeFiles: composeFiles,
			ProjectName:  taskName,
			Services:     []string{service},
		}

		taskDefinitionInput, err := t.Transform()
		if err != nil {
			return err
		}

		clusterOutput, err := api.StackOutputs(svc.Cloudformation, *clusterStack.StackName)
		if err != nil {
			return err
		}

		logGroup, exists := clusterOutput["LogGroupName"]
		if !exists {
			return fmt.Errorf("Expected to find a LogGroupName in stack output")
		}

		log.Printf("Setting tasks to use log group %s", logGroup)
		for _, def := range taskDefinitionInput.ContainerDefinitions {
			if def.LogConfiguration == nil {
				def.LogConfiguration = &ecs.LogConfiguration{
					LogDriver: aws.String("awslogs"),
					Options: map[string]*string{
						"awslogs-group":         aws.String(logGroup),
						"awslogs-region":        aws.String(os.Getenv("AWS_REGION")),
						"awslogs-stream-prefix": aws.String(taskName),
					},
				}
			}
		}

		log.Printf("Registering a task for %s", taskName)
		resp, err := svc.ECS.RegisterTaskDefinition(taskDefinitionInput)
		if err != nil {
			return err
		}

		taskDefinition := fmt.Sprintf("%s:%d",
			*resp.TaskDefinition.Family, *resp.TaskDefinition.Revision)

		runTaskInput := &ecs.RunTaskInput{
			TaskDefinition: aws.String(taskDefinition),
			Cluster:        aws.String(cluster),
			Count:          aws.Int64(1),
			Overrides: &ecs.TaskOverride{
				ContainerOverrides: []*ecs.ContainerOverride{},
			},
		}

		if len(commands) > 0 {
			cmds := []*string{}

			for _, command := range commands {
				cmds = append(cmds, aws.String(command))
			}

			runTaskInput.Overrides.ContainerOverrides = append(
				runTaskInput.Overrides.ContainerOverrides,
				&ecs.ContainerOverride{
					Command: cmds,
					Name:    aws.String(service),
				},
			)
		}

		log.Printf("Running task %s", taskDefinition)
		runResp, err := svc.ECS.RunTask(runTaskInput)
		if err != nil {
			return err
		}

		taskARNs := []*string{}
		for _, t := range runResp.Tasks {
			taskARNs = append(taskARNs, t.TaskArn)
		}

		go func() {
			err = svc.ECS.WaitUntilTasksStopped(&ecs.DescribeTasksInput{
				Cluster: aws.String(cluster),
				Tasks:   taskARNs,
			})
			if err != nil {
				log.Fatal(err)
			}

			output, err := svc.ECS.DescribeTasks(&ecs.DescribeTasksInput{
				Cluster: aws.String(cluster),
				Tasks:   taskARNs,
			})
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Task finished: %#v", output)
			var exitCode int

			for _, task := range output.Tasks {
				for _, container := range task.Containers {
					if *container.ExitCode != 0 {
						exitCode = int(*container.ExitCode)
						log.Printf("Container %s exited with %d", *container.Name, exitCode)
					}
				}
			}

			// FIX: gross, but logs lag behind
			time.Sleep(time.Second * 5)
			os.Exit(exitCode)
		}()

		containerID := path.Base(*runResp.Tasks[0].TaskArn)
		prefix := fmt.Sprintf("%s/%s/%s", taskName, service, containerID)

		log.Printf("Following logs for %s", prefix)

		w := &logWatcher{
			LogGroup:  logGroup,
			LogPrefix: prefix,
			services:  svc,
			Printer: func(ev *logs.FilteredLogEvent) {
				fmt.Println(*ev.Message)
			},
		}

		if err := w.Watch(context.Background()); err != nil {
			return err
		}

		return nil
	})
}
