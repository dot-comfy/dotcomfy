/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cobra

import (
	"os"
	"os/user"

	"github.com/spf13/cobra"

	Log "dotcomfy/internal/logger"
	"dotcomfy/internal/services"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Pushes your local changes",
	Long: `Pushes your local changes to the remote repository on the current
branch. Note that you must have write permissions for this to succeed.`,
	Run: func(cmd *cobra.Command, args []string) {
		LOGGER := Log.GetLogger()
		user, err := user.Current()
		if err != nil {
			LOGGER.Fatal(err)
		}
		dotcomfy_dir := user.HomeDir + "/.dotcomfy"
		err = services.Push(dotcomfy_dir)
		if err != nil {
			LOGGER.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
