/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cobra

import (
	"fmt"
	"os"
	"os/user"

	"github.com/spf13/cobra"

	Log "dotcomfy/internal/logger"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncs your dotcomfy installation with the remote repository",
	Long: `This command will pull the latest changes from the remote repository
and update your dotcomfy installation accordingly.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sync called")
	},
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Pushes your local changes",
	Long: `Pushes your local changes to the remote repository on the current
branch. Note that you must have write permissions for this to succeed.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("push called")

		LOGGER = Log.GetLogger()
		user, err := user.Current()
		if err != nil {
			LOGGER.Fatal(err)
		}
		dotcomfy_dir := user.HomeDir + "/.dotcomfy"
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
