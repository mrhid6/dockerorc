package node

import (
	"context"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/mrhid6/dockerorc/services/dockerorc-adm/api"
	"github.com/mrhid6/dockerorc/utils"
	"github.com/spf13/cobra"
)

var clusterNodeAddMasterIP string

func init() {
	Cmd.AddCommand(clusterNodeAddCmd)
}

var clusterNodeAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a node to the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		nodeName, err := os.Hostname()
		if err != nil {
			panic(err)
		}

		if err := RegisterNode(nodeName); err != nil {
			panic(err)
		}

		CreateNodeContainer(nodeName)
	},
}

func init() {
	clusterNodeAddCmd.Flags().StringVarP(&clusterNodeAddMasterIP, "masterip", "m", "", "The Master IP address")
	clusterNodeAddCmd.MarkFlagRequired("masterip")
}

func CheckNodeContainer(nodeName string) bool {
	containerName := "dockerorc-node-" + nodeName
	container, err := utils.GetDockerContainerByName(containerName)

	if err != nil {
		return false
	}

	return container != nil
}

func CreateNodeContainer(nodeName string) error {
	log.Println("Creating node container...")

	if CheckNodeContainer(nodeName) {
		log.Println("Node Container already exists!")
		return nil
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	containerName := "dockerorc-node-" + nodeName

	res, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:    "dockerorc-node:latest",
			Hostname: containerName,
			Env: []string{
				"DOCKERORC_MASTERIP=" + clusterNodeAddMasterIP,
			},
		},
		&container.HostConfig{
			NetworkMode: "orcbridge",
			Binds:       []string{"/var/run/docker.sock:/var/run/docker.sock"},
		},
		&network.NetworkingConfig{},
		nil,
		containerName,
	)

	if err != nil {
		return err
	}

	err = cli.ContainerStart(context.Background(), res.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	log.Println("Successfully created node container!")

	return nil
}

func RegisterNode(nodeName string) error {
	log.Println("Registering node with cluster...")
	type RegisterNodeBody struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
	}

	body := RegisterNodeBody{
		Name: nodeName,
		IP:   utils.GetOutboundIP().String(),
	}

	var res interface{}

	var url = "http://:6443/api/nodes/register"

	if err := api.SendPostRequest(url, body, &res); err != nil {
		return err
	}

	log.Println("Registered node successfully!")

	return nil
}
