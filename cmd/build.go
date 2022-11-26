/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/aura-studio/dynamic/builder"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("warehouse").Value.String() != "" {
			builder.DefaultConfig.WareHouse = cmd.Flag("warehouse").Value.String()
		}
		if cmd.Flag("gover").Value.String() != "" {
			builder.DefaultConfig.GoVer = cmd.Flag("gover").Value.String()
		}
		if len(args) > 0 {
			if strings.Contains(args[0], "@") {
				builder.BuildFromRemote(args[0], args[1:]...)
			} else {
				builder.BuildFromJSONPath(args[0])
			}
		}
		if cmd.Flag("config").Value.String() != "" {
			builder.BuildFromJSONFile(cmd.Flag("config").Value.String())
		}
		if cmd.Flag("path").Value.String() != "" {
			builder.BuildFromJSONPath(cmd.Flag("path").Value.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	buildCmd.Flags().StringP("config", "c", "", "path of config file")
	buildCmd.Flags().StringP("path", "p", "", "path of work directory")
	buildCmd.Flags().StringP("warehouse", "w", "", "path of warehouse")
	buildCmd.Flags().StringP("gover", "g", "", "path of go version")
}
