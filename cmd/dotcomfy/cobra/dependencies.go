/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cobra

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"dotcomfy/internal/config"
)

var all bool

// dependenciesCmd represents the dependencies command
var dependenciesCmd = &cobra.Command{
	Use:   "dependencies",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if all {
			allDependencies()
			os.Exit(0)
		}

		dependency := args[0]

		dependency_map, err := config.GetDependency(dependency)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for k, v := range dependency_map {
			switch k {
			case "version":
				fmt.Printf("Version: %s\n", v)
			case "steps":
				for i, step := range v.([]interface{}) {
					fmt.Printf("Step %d: %s\n", i, step)
				}
			case "post_installation_steps":
			}
		}
	},
}

func allDependencies() {
	dependencies := config.GetDependencies()
	fmt.Println(dependencies)
}

func init() {
	rootCmd.AddCommand(dependenciesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dependenciesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dependenciesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	dependenciesCmd.Flags().BoolVar(&all, "all", false, "Get all dependencies")
}
