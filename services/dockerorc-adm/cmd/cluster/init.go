package cluster

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/mrhid6/dockerorc/services/dockerorc-adm/api"
	"github.com/mrhid6/dockerorc/utils"
	"github.com/spf13/cobra"
)

var clusterInitNetwork string

func init() {
	Cmd.AddCommand(clusterInitCmd)
}

var clusterInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialises cluster",
	Run: func(cmd *cobra.Command, args []string) {
		CreateNetwork()
		err := CreateSystemContainer()
		if err != nil {
			panic(err)
		}

		CheckSystemIsRunning()

		ipAddress := utils.GetOutboundIP()

		log.Println("Cluster has successfully been set up.")
		fmt.Printf("Cluster Master IP Address: %s:6443\r\n", ipAddress.String())
	},
}

func init() {
	clusterInitCmd.Flags().StringVarP(&clusterInitNetwork, "network", "n", "", "Container Network")
	clusterInitCmd.MarkFlagRequired("network")
}

func CreateNetwork() error {

	log.Println("Creating Cluster Network...")
	if CheckNetwork() {
		log.Println("Cluster Network has already been created. Skipping.")
		return nil
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	newnetwork := types.NetworkCreate{IPAM: &network.IPAM{
		Driver: "default",
		Config: []network.IPAMConfig{network.IPAMConfig{
			Subnet:  clusterInitNetwork,
			IPRange: clusterInitNetwork,
		}},
	}}

	_, err = cli.NetworkCreate(context.Background(), "orcbridge", newnetwork)
	if err != nil {
		return err
	}

	return nil
}

func CheckNetwork() bool {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	_, err = cli.NetworkInspect(context.Background(), "orcbridge", types.NetworkInspectOptions{})
	if err != nil {
		return false
	}
	return true
}

func CreateSystemContainer() error {

	log.Println("Creating Docker Orc System Container..")
	if CheckSystemContainer() {
		log.Println("Skipped docker container creation - already exists..")
		return nil
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	var portBindings = nat.PortMap{
		"6443/tcp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "6443",
			},
		},
	}
	var exposedPorts = nat.PortSet{
		"6443/tcp": struct{}{},
	}

	res, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        "dockerorc-system:latest",
			Hostname:     "dockerorc-system",
			ExposedPorts: exposedPorts,
		},
		&container.HostConfig{
			NetworkMode:  "orcbridge",
			PortBindings: portBindings,
		},
		&network.NetworkingConfig{},
		nil,
		"dockerorc-system",
	)

	if err != nil {
		return err
	}

	err = cli.ContainerStart(context.Background(), res.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return nil
}

func CheckSystemContainer() bool {
	container, err := utils.GetDockerContainerByName("dockerorc-system")

	if err != nil {
		return false
	}

	return container != nil
}

func CheckSystemIsRunning() {
	log.Println("Checking connection to system container")
	var res interface{}

	timeoutCounter := 0

	for {
		err := api.SendGetRequest("http://:6443/api/", &res)

		if timeoutCounter >= 5 {
			log.Println("System container timed out.")
		}

		if err != nil {
			log.Println("System still not ready. Retrying in 5 seconds")
			timeoutCounter++
		} else {
			break
		}
		time.Sleep(5 * time.Second)
	}
	log.Println("System Container is ready!")
}
