package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mrhid6/dockerorc/services/dockerorc-node/api"
)

func main() {
	wait := gracefulShutdown(context.Background(), 30*time.Second, map[string]operation{
		"main": func(ctx context.Context) error {
			return ShutdownMain()
		},
	})

	SetNodeStatus(true)

	<-wait
}

func ShutdownMain() error {
	if err := SetNodeStatus(false); err != nil {
		return err
	}
	return nil
}

func GetContainers() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})

	if err != nil {
		return err
	}

	for _, container := range containers {
		fmt.Println(container.Names[0])
	}

	return nil
}

func SetNodeStatus(status bool) error {

	var masterIp = os.Getenv("DOCKERORC_MASTERIP")
	var url = "http://" + masterIp + ":6443/api/nodes/status"

	type NodeStatusBody struct {
		Name   string `json:"name"`
		Status bool   `json:"status"`
	}

	hostname, _ := os.Hostname()
	nodeName := strings.Replace(hostname, "dockerorc-node-", "", -1)

	body := NodeStatusBody{
		Name:   nodeName,
		Status: status,
	}

	var res interface{}

	err := api.SendPostRequest(url, body, &res)
	if err != nil {
		return err
	}
	return nil
}
