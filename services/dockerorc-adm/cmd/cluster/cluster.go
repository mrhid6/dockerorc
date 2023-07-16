package cluster

import (
	"github.com/mrhid6/dockerorc/services/dockerorc-adm/cmd/cluster/node"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "cluster",
	Short: "Docker Orc Admin Cluster",
}

func init() {
	Cmd.AddCommand(node.Cmd)
}
