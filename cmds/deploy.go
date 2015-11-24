package cmds

import (
	"path/filepath"
	"strings"

	"github.com/99designs/ecs-former/core"
)

type DeployCommandInput struct {
	TaskFile  string
	TaskName  string
	Templates core.Templates
}

func DeployCommand(ui *Ui, input DeployCommandInput) {
	if input.TaskName == "" {
		input.TaskName = strings.Split(filepath.Base(input.TaskFile), ".")[0]
	}
}
