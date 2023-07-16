package cmd

import (
	"fmt"
	"os"

	"github.com/mrhid6/dockerorc/services/dockerorc-adm/cmd/cluster"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dockerorc-adm",
	Short: "Docker Orc Admin",
	Long:  "Docker Orc Administration commands",
}

func init() {
	rootCmd.AddCommand(cluster.Cmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
