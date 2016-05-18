package compose

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
)

func TransformComposeFile(composeFile string, projectName string) (*ecs.RegisterTaskDefinitionInput, error) {
	task := ecs.RegisterTaskDefinitionInput{
		Family:               aws.String(projectName),
		ContainerDefinitions: []*ecs.ContainerDefinition{},
		Volumes:              []*ecs.Volume{},
	}

	project, err := docker.NewProject(&docker.Context{
		Context: project.Context{
			ComposeFiles: []string{composeFile},
			ProjectName:  projectName,
		}})
	if err != nil {
		return nil, err
	}


	for _, name := range project.Configs.Keys() {
		config, _ := project.Configs.Get(name)
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
			def.Cpu = aws.Int64(config.CPUShares)
		}

		if config.MemLimit > 0 {
			def.Memory = aws.Int64(config.MemLimit)
		}

		if config.Privileged {
			def.Privileged = aws.Bool(config.Privileged)
		}

		if slice := config.DNS.Slice(); len(slice) > 0 {
			for _, dns := range slice {
				def.DnsServers = append(def.DnsServers, aws.String(dns))
			}
		}

		if slice := config.DNSSearch.Slice(); len(slice) > 0 {
			for _, item := range slice {
				def.DnsSearchDomains = append(def.DnsSearchDomains, aws.String(item))
			}
		}

		if cmds := config.Command.Slice(); len(cmds) > 0 {
			def.Command = []*string{}
			for _, command := range cmds {
				def.Command = append(def.Command, aws.String(command))
			}
		}

		if cmds := config.Entrypoint.Slice(); len(cmds) > 0 {
			def.EntryPoint = []*string{}
			for _, command := range cmds {
				def.EntryPoint = append(def.EntryPoint, aws.String(command))
			}
		}

		if slice := config.Environment.Slice(); len(slice) > 0 {
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

		if links := config.Links.Slice(); len(links) > 0 {
			def.Links = []*string{}
			for _, link := range links {
				def.Links = append(def.Links, aws.String(link))
			}
		}

		if vols := config.Volumes; len(vols) > 0 {
			def.MountPoints = []*ecs.MountPoint{}
			for idx, vol := range vols {
				parts := strings.Split(vol, ":")
				volumeName := fmt.Sprintf("%s-vol%d", name, idx)

				volume := ecs.Volume{
					Host: &ecs.HostVolumeProperties{
						SourcePath: aws.String(parts[0]),
					},
					Name: aws.String(volumeName),
				}

				mount := ecs.MountPoint{
					SourceVolume:  aws.String(volumeName),
					ContainerPath: aws.String(parts[1]),
				}

				if len(parts) == 3 && parts[2] == "ro" {
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

		if config.Build != "" {
			return nil, errors.New("Build directive not supported")
		}
		if config.Dockerfile != "" {
			return nil, errors.New("Dockerfile directive not supported")
		}
		if config.DomainName != "" {
			return nil, errors.New("DomainName directive not supported")
		}
		if config.VolumeDriver != "" {
			return nil, errors.New("VolumeDriver directive not supported")
		}

		task.ContainerDefinitions = append(task.ContainerDefinitions, &def)

		// Not Implemented
		// CapAdd:[]string(nil),
		// CapDrop:[]string(nil),
		// CPUSet:"",
		// ContainerName:"",
		// Devices:[]string(nil),
		// EnvFile:project.Stringorslice{parts:[]string(nil)},
		// Labels:project.SliceorMap{parts:map[string]string(nil)},
		// LogDriver:"",
		// MemSwapLimit:0,
		// Net:"",
		// Pid:"",
		// Uts:"",
		// Ipc:"",
		// Restart:"",
		// ReadOnly:false,
		// StdinOpen:false,
		// SecurityOpt:[]string(nil),
		// Tty:false,
		// User:"",
		// Expose:[]string(nil),
		// ExternalLinks:[]string(nil),
		// LogOpt:map[string]string(nil),
		// ExtraHosts:[]string(nil)}
	}

	return &task, nil
}
