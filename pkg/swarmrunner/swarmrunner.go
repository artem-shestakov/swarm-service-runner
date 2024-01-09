package swarmrunner

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"gopkg.in/yaml.v3"
)

func RunService(file string) {
	var yamlStruct map[string]interface{}

	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = yaml.Unmarshal(data, &yamlStruct)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	for svcName, svcSpec := range yamlStruct["services"].(map[string]interface{}) {
		fmt.Println(svcName)
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

		fmt.Println(svc)
		svcResp, err := cli.ServiceCreate(ctx, svc, types.ServiceCreateOptions{})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(svcResp)

	}
}
