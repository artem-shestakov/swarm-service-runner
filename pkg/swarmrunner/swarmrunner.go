package swarmrunner

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func CreateServices(file string) (*[]swarm.ServiceSpec, error) {
	var services []swarm.ServiceSpec
	var yamlStruct map[string]interface{}

	logrus.Debugf("Reading file %s", file)
	data, err := os.ReadFile(file)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, err
	}

	logrus.Debugf("Parsing yaml struct")
	err = yaml.Unmarshal(data, &yamlStruct)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, err
	}

	for svcName, svcSpec := range yamlStruct["services"].(map[string]interface{}) {
		swarmPorts := []swarm.PortConfig{}
		if svcSpec.(map[string]interface{})["ports"] != nil {
			for _, pubPorts := range svcSpec.(map[string]interface{})["ports"].([]interface{}) {
				var pubPort, targetPort int
				if l := strings.Split(pubPorts.(string), ":"); len(l) > 0 {
					pubPort, _ = strconv.Atoi(strings.Split(pubPorts.(string), ":")[0])
					targetPort, _ = strconv.Atoi(strings.Split(pubPorts.(string), ":")[1])
				} else {
					pubPort, _ = strconv.Atoi(strings.Split(pubPorts.(string), ":")[0])
					targetPort = pubPort
				}

				swarmPorts = append(swarmPorts, swarm.PortConfig{
					PublishedPort: uint32(pubPort),
					TargetPort:    uint32(targetPort),
				})
			}
		}

		svc := swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: svcName,
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: &swarm.ContainerSpec{
					Image: svcSpec.(map[string]interface{})["image"].(string),
				},
			},
			EndpointSpec: &swarm.EndpointSpec{
				Mode:  swarm.ResolutionModeVIP,
				Ports: swarmPorts,
			},
		}
		services = append(services, svc)
	}
	return &services, nil
}

func CreateService(service swarm.ServiceSpec) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Panicln(err)
	}
	defer cli.Close()

	if ok, err := isServiceExists(ctx, cli, service.Name); err != nil {
		return err
	} else if ok {
		logrus.Warnf("Service %s already exists", service.Name)
	} else {
		logrus.Infof("Creating service %s", service.Name)
		svcResp, err := cli.ServiceCreate(ctx, service, types.ServiceCreateOptions{})
		if err != nil {
			logrus.Errorln(err.Error())
			return err
		}
		if svcResp.ID != "" {
			logrus.Infof("Service %s created with ID %s", service.Name, svcResp.ID)
		}
		for _, svcWarn := range svcResp.Warnings {
			logrus.Warn(svcWarn)
		}
	}
	return nil
}

func isServiceExists(ctx context.Context, cli *client.Client, svcName string) (bool, error) {
	svcList, err := cli.ServiceList(ctx, types.ServiceListOptions{})
	if err != nil {
		return false, err
	}
	for _, svc := range svcList {
		if svc.Spec.Name == svcName {
			return true, nil
		}
	}
	return false, nil
}
