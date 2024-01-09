/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "swarm-runner",
	Short: "Create Docker Swarm services from Compose file without Swarm stack",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logLevel, _ := cmd.Flags().GetString("verbosity")
		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			logrus.Panicln(err.Error())
		}
		logrus.SetLevel(level)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("verbosity", "v", logrus.InfoLevel.String(), "Set logger level. Use: debug, info, warn, error, fatal or panic")
}
