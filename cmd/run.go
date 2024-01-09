/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

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
			logrus.Debugln(err.Error())
		}
		services, err := swarmrunner.CreateServices(file)
		if err != nil {
			logrus.Errorln(err.Error())
			os.Exit(1)
		}
		for _, svc := range *services {
			swarmrunner.CreateService(svc)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("file", "f", "docker-compose.yaml", "Path to docker compose file")
}
