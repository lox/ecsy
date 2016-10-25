package main

import (
	"os"

	"github.com/lox/ecsy/api"
	"github.com/lox/ecsy/cmd"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	Version string
)

func main() {
	run(os.Args[1:], os.Exit)
}

func run(args []string, exit func(code int)) {
	app := kingpin.New("ecsy",
		`A tool for managing and deploying ECS clusters`)

	app.Version(Version)
	app.Writer(os.Stdout)
	app.DefaultEnvars()
	app.Terminate(exit)

	cmd.ConfigureCreateCluster(app, api.DefaultServices)
	cmd.ConfigureDeleteCluster(app, api.DefaultServices)
	cmd.ConfigureCreateService(app, api.DefaultServices)
	cmd.ConfigurePoll(app, api.DefaultServices)
	cmd.ConfigureDeploy(app, api.DefaultServices)
	cmd.ConfigureDumpTaskDefinition(app, api.DefaultServices)

	kingpin.MustParse(app.Parse(args))
}
