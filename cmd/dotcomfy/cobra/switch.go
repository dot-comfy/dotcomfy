/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cobra

import (
	"fmt"

	"github.com/spf13/cobra"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch [OPTION]... NAME...",
	Short: "Switch to a different branch or repository of dotfiles",
	Long: `Switch to a different branch or repository of dotfiles.
	This will uninstall your current dotcomfy installation,
	and install the dotfiles from the branch or repository
	you are switching to.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("switch called")
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// switchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// switchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	switchCmd.Flags().StringVar(&branch, "branch", "", "Branch to switch to")
	switchCmd.Flags().StringVar(&repo, "repo", "", "Repository to switch to")
}
