package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	logs "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/lox/ecsy/api"
	"golang.org/x/net/context"
	"gopkg.in/alecthomas/kingpin.v2"
)

func parseEventTime(t int64) time.Time {
	return time.Unix(0, t*int64(time.Millisecond))
}

func ConfigureLogs(app *kingpin.Application, svc api.Services) {
	var logGroup, prefix string

	cmd := app.Command("logs", "Follow a cloudwatch log")
	cmd.Flag("log-group", "The cloudwatch logs group to follow").
		Short('g').
		Required().
		StringVar(&logGroup)

	cmd.Flag("prefix", "Only show logs with this prefix").
		Short('p').
		StringVar(&prefix)

	cmd.Action(func(c *kingpin.ParseContext) error {
		w := &logWatcher{
			LogGroup:  logGroup,
			LogPrefix: prefix,
			services:  svc,
		}

		if err := w.Watch(context.Background()); err != nil {
			return err
		}

		return nil
	})
}

type logWatcher struct {
	LogGroup, LogPrefix string
	Printer             func(ev *logs.FilteredLogEvent)
	services            api.Services
}

func (lw *logWatcher) findStreams() ([]*string, error) {
	params := &logs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(lw.LogGroup),
		LogStreamNamePrefix: aws.String(lw.LogPrefix),
		Descending:          aws.Bool(true),
	}

	streams := []*string{}
	err := lw.services.Logs.DescribeLogStreamsPages(params, func(page *logs.DescribeLogStreamsOutput, lastPage bool) bool {
		for _, stream := range page.LogStreams {
			streams = append(streams, stream.LogStreamName)
		}
		return lastPage
	})

	return streams, err
}

func (lw *logWatcher) printEventsAfter(ts int64) (int64, error) {
	var streams []*string
	var err error

	// sometimes the stream takes a while to appear
	for i := 0; i < 10; i++ {
		streams, err = lw.findStreams()
		if err != nil || len(streams) == 0 {
			time.Sleep(time.Second * 2)
		}
	}

	if err != nil {
		return ts, err
	} else if len(streams) == 0 {
		return ts, errors.New("No stream found")
	}

	filterInput := &logs.FilterLogEventsInput{
		LogGroupName:   aws.String(lw.LogGroup),
		LogStreamNames: streams,
		StartTime:      aws.Int64(ts + 1),
	}

	err = lw.services.Logs.FilterLogEventsPages(filterInput,
		func(p *logs.FilterLogEventsOutput, lastPage bool) (shouldContinue bool) {
			for _, event := range p.Events {
				ts = *event.Timestamp
				if lw.Printer == nil {
					lw.printEvent(event)
				} else {
					lw.Printer(event)
				}
			}
			return lastPage
		})

	return ts, err
}

func (lw *logWatcher) printEvent(ev *logs.FilteredLogEvent) {
	name := strings.TrimPrefix(*ev.LogStreamName, lw.LogPrefix)
	if len(name) > 50 {
		name = name[:47] + "..."
	}

	fmt.Printf("%-20s %-52s %s\n",
		parseEventTime(*ev.Timestamp).Format(time.Stamp),
		name,
		*ev.Message,
	)
}

func (lw *logWatcher) Watch(ctx context.Context) error {
	var after int64
	var err error

	for {
		select {
		case <-time.After(1 * time.Second):
			if after, err = lw.printEventsAfter(after); err != nil {
				return err
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
