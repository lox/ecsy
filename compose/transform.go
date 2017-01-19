package compose

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/docker/libcompose/config"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
)

type Transformer struct {
	ComposeFiles      []string
	ProjectName       string
	Services          []string
	EnvironmentLookup config.EnvironmentLookup
}

func (t *Transformer) Transform() (*ecs.RegisterTaskDefinitionInput, error) {
	task := ecs.RegisterTaskDefinitionInput{
		Family:               aws.String(t.ProjectName),
		ContainerDefinitions: []*ecs.ContainerDefinition{},
		Volumes:              []*ecs.Volume{},
	}

	projectCtx := project.Context{
		ComposeFiles: t.ComposeFiles,
		ProjectName:  t.ProjectName,
	}

	if t.EnvironmentLookup != nil {
		projectCtx.EnvironmentLookup = t.EnvironmentLookup
	}

	p, err := docker.NewProject(&ctx.Context{Context: projectCtx}, nil)
	if err != nil {
		return nil, err
	}

	for _, name := range p.(*project.Project).ServiceConfigs.Keys() {
		if !isServiceIncluded(name, t.Services) {
			log.Printf("Skipping service %s", name)
			continue
		}

		config, _ := p.GetServiceConfig(name)
		var def = ecs.ContainerDefinition{
			Name: aws.String(name),
		}

		if config.Image != "" {
			def.Image = aws.String(config.Image)
		}

		if config.Hostname != "" {
			def.Hostname = aws.String(config.Hostname)
		}

		if config.WorkingDir != "" {
			def.WorkingDirectory = aws.String(config.WorkingDir)
		}

		if config.CPUShares > 0 {
			def.Cpu = aws.Int64(int64(config.CPUShares))
		}

		// docker-compose expresses it's limit in bytes, ECS in mb
		if config.MemLimit > 0 {
			def.Memory = aws.Int64(int64(config.MemLimit / 1024 / 1024))
		}

		if config.MemReservation > 0 {
			def.MemoryReservation = aws.Int64(int64(config.MemReservation / 1024 / 1024))
		}

		if config.Privileged {
			def.Privileged = aws.Bool(config.Privileged)
		}

		if slice := []string(config.DNS); len(slice) > 0 {
			for _, dns := range slice {
				def.DnsServers = append(def.DnsServers, aws.String(dns))
			}
		}

		if slice := []string(config.DNSSearch); len(slice) > 0 {
			for _, item := range slice {
				def.DnsSearchDomains = append(def.DnsSearchDomains, aws.String(item))
			}
		}

		if cmds := []string(config.Command); len(cmds) > 0 {
			def.Command = []*string{}
			for _, command := range cmds {
				def.Command = append(def.Command, aws.String(command))
			}
		}

		if cmds := []string(config.Entrypoint); len(cmds) > 0 {
			def.EntryPoint = []*string{}
			for _, command := range cmds {
				def.EntryPoint = append(def.EntryPoint, aws.String(command))
			}
		}

		if slice := []string(config.Environment); len(slice) > 0 {
			def.Environment = []*ecs.KeyValuePair{}
			for _, val := range slice {
				parts := strings.SplitN(val, "=", 2)
				def.Environment = append(def.Environment, &ecs.KeyValuePair{
					Name:  aws.String(parts[0]),
					Value: aws.String(parts[1]),
				})
			}
		}

		if ports := config.Ports; len(ports) > 0 {
			def.PortMappings = []*ecs.PortMapping{}
			for _, val := range ports {
				parts := strings.Split(val, ":")
				mapping := &ecs.PortMapping{}

				// TODO: support host to map to
				if len(parts) > 0 {
					portInt, err := strconv.ParseInt(parts[0], 10, 64)
					if err != nil {
						return nil, err
					}
					mapping.ContainerPort = aws.Int64(portInt)
				}

				if len(parts) > 1 {
					hostParts := strings.Split(parts[1], "/")
					portInt, err := strconv.ParseInt(hostParts[0], 10, 64)
					if err != nil {
						return nil, err
					}
					mapping.HostPort = aws.Int64(portInt)

					// handle the protocol at the end of the mapping
					if len(hostParts) > 1 {
						mapping.Protocol = aws.String(hostParts[1])
					}
				}

				if len(parts) == 0 || len(parts) > 2 {
					return nil, errors.New("Unsupported port mapping " + val)
				}

				def.PortMappings = append(def.PortMappings, mapping)
			}
		}

		if links := []string(config.Links); len(links) > 0 {
			def.Links = []*string{}
			for _, link := range links {
				def.Links = append(def.Links, aws.String(link))
			}
		}

		if dependsOn := []string(config.DependsOn); len(dependsOn) > 0 {
			def.Links = []*string{}
			for _, link := range dependsOn {
				def.Links = append(def.Links, aws.String(link))
			}
		}

		if config.Volumes != nil {
			def.MountPoints = []*ecs.MountPoint{}
			for idx, vol := range config.Volumes.Volumes {
				volumeName := fmt.Sprintf("%s-vol%d", name, idx)

				volume := ecs.Volume{
					Host: &ecs.HostVolumeProperties{
						SourcePath: aws.String(vol.Source),
					},
					Name: aws.String(volumeName),
				}

				mount := ecs.MountPoint{
					SourceVolume:  aws.String(volumeName),
					ContainerPath: aws.String(vol.Destination),
				}

				if vol.AccessMode == "ro" {
					mount.ReadOnly = aws.Bool(true)
				}

				task.Volumes = append(task.Volumes, &volume)
				def.MountPoints = append(def.MountPoints, &mount)
			}
		}

		if volsFrom := config.VolumesFrom; len(volsFrom) > 0 {
			def.VolumesFrom = []*ecs.VolumeFrom{}
			for _, container := range volsFrom {
				def.VolumesFrom = append(def.VolumesFrom, &ecs.VolumeFrom{
					SourceContainer: aws.String(container),
				})
			}
		}

		if config.Logging.Driver != "" {
			def.LogConfiguration = &ecs.LogConfiguration{
				LogDriver: aws.String(config.Logging.Driver),
				Options:   map[string]*string{},
			}
			for k, v := range config.Logging.Options {
				def.LogConfiguration.Options[k] = aws.String(v)
			}
		}

		// TODO: the following directives are ignored / not handled
		// Networks
		// Expose

		for _, i := range []struct {
			Key   string
			Value []string
		}{
			{"CapAdd", config.CapAdd},
			{"CapDrop", config.CapDrop},
			{"Devices", config.Devices},
			{"EnvFile", config.EnvFile},
			{"SecurityOpt", config.SecurityOpt},
			{"ExternalLinks", config.ExternalLinks},
			{"ExtraHosts", config.ExtraHosts},
			{"Extends", config.Extends},
			{"GroupAdd", config.GroupAdd},
			{"DNSOpts", config.DNSOpts},
			{"Tmpfs", config.Tmpfs},
		} {
			if len(i.Value) > 0 {
				return nil, fmt.Errorf("%s directive not supported", i.Key)
			}
		}

		for _, i := range []struct {
			Key   string
			Value string
		}{
			{"Build.Context", config.Build.Context},
			{"Build.Dockerfile", config.Build.Dockerfile},
			{"DomainName", config.DomainName},
			{"VolumeDriver", config.VolumeDriver},
			{"CPUSet", config.CPUSet},
			{"NetworkMode", config.NetworkMode},
			{"Pid", config.Pid},
			{"Uts", config.Uts},
			{"Ipc", config.Ipc},
			{"Restart", config.Restart},
			{"User", config.User},
			{"StopSignal", config.StopSignal},
			{"MacAddress", config.MacAddress},
			{"Isolation", config.Isolation},
			{"NetworkMode", config.NetworkMode},
			{"CgroupParent", config.CgroupParent},
		} {
			if i.Value != "" {
				return nil, fmt.Errorf("%s directive not supported", i.Key)
			}
		}

		for _, i := range []struct {
			Key   string
			Value int
		}{
			{"CPUShares", int(config.CPUShares)},
			{"CPUQuota", int(config.CPUQuota)},
			{"MemSwapLimit", int(config.MemSwapLimit)},
			{"MemSwappiness", int(config.MemSwappiness)},
			{"OomScoreAdj", int(config.OomScoreAdj)},
			{"ShmSize", int(config.ShmSize)},
		} {
			if i.Value != 0 {
				return nil, fmt.Errorf("%s directive not supported", i.Key)
			}
		}

		for _, i := range []struct {
			Key   string
			Value bool
		}{
			{"ReadOnly", config.ReadOnly},
			{"StdinOpen", config.StdinOpen},
			{"Tty", config.Tty},
		} {
			if i.Value {
				return nil, fmt.Errorf("%s directive not supported", i.Key)
			}
		}

		if config.Labels != nil {
			return nil, fmt.Errorf("Labels directive not supported")
		}

		if len(config.Ulimits.Elements) > 0 {
			return nil, fmt.Errorf("Ulimits directive not supported")
		}

		task.ContainerDefinitions = append(task.ContainerDefinitions, &def)
	}

	return &task, nil
}

func isServiceIncluded(name string, included []string) bool {
	if len(included) == 0 {
		return true
	}
	for _, val := range included {
		if name == val {
			return true
		}
	}
	return false
}

type envMap map[string][]string

func (em envMap) Lookup(key string, config *config.ServiceConfig) []string {
	results := []string{}
	for _, val := range em[key] {
		results = append(results, key+"="+val)
	}
	return results
}
