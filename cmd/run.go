/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/artem-shestakov/swarm-service-runner/pkg/swarmrunner"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Docker swarm services",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			logrus.Debugf("Couldn't get flag 'file'. %s", err.Error())
		}
		swarmrunner.RunService(file)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("file", "f", "docker-compose.yaml", "Path to docker compose file")
}
