/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cobra

import (
	"dotcomfy/internal/config"
	"dotcomfy/internal/services"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [GitHub username/repo URL]",
	Short: "Install dotfiles from a Git repo",
	Long: `Install dotfiles from a Git repo. You can pass in just a GitHub username
	(which will look for the repository "https://github.com/{username}/dotfiles.git"),
	or the full URL to a Git repo containing dotfiles`,
	Args: cobra.MinimumNArgs(1),
	Run:  run,
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func run(cmd *cobra.Command, args []string) {
	path := config.GetConfig().Foo

	if path == "" {
		tmp, err := os.MkdirTemp("", "dotcomfy-")
		fmt.Println("install called")
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %+v\n", err.Error())
			os.Exit(1)
		}
		defer func() {
			os.RemoveAll(tmp)
		}()
		path = tmp
	}

	if err := services.CloneTo(args[0], path); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %+v\n", err)
		os.Exit(1)
	}
}
