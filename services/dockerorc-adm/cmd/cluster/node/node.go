package node

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "nodes",
	Short: "Docker Orc Admin Cluster Nodes",
}
