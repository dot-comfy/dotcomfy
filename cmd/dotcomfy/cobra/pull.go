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

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Syncs your dotcomfy installation with the remote repository",
	Long: `This command will pull the latest changes from the remote repository
and update your dotcomfy installation accordingly.`,
	Run: func(cmd *cobra.Command, args []string) {
		LOGGER := Log.GetLogger()
		user, err := user.Current()
		if err != nil {
			LOGGER.Fatal(err)
		}
		dotcomfy_dir := user.HomeDir + "/.dotcomfy"

		err = services.Pull(dotcomfy_dir)
		if err != nil {
			LOGGER.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
