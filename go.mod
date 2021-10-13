module github.com/lox/ecsy

go 1.16

replace github.com/docker/docker => github.com/docker/engine v0.0.0-20190725163905-fa8dd90ceb7b

require (
	github.com/aws/aws-sdk-go v1.5.6-0.20161115230331-55795bd6e929
	github.com/docker/libcompose v0.4.1-0.20210616120443-2a046c0bdbf2
	github.com/fatih/color v1.1.1-0.20161228204310-9ab0325f4904
	github.com/go-ini/ini v1.63.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattn/go-colorable v0.0.7 // indirect
	github.com/mattn/go-isatty v0.0.0-20161123143637-30a891c33c7c // indirect
	golang.org/x/net v0.0.0-20190613194153-d28f0bde5980
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)
