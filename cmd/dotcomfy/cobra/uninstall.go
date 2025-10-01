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
	"dotcomfy/internal/services"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Removes current dotcomfy installation",
	Long: `Removes current dotcomfy installation
	from your system and restores previous dotfiles.`,
	Run: func(cmd *cobra.Command, args []string) {
		LOGGER = Log.GetLogger()
		var confirmation string
		if !CONFIRM {
			LOGGER.Info("Are you sure you want to uninstall the current dotcomfy installation? (y/n)")
			fmt.Print("Are you sure you want to uninstall the current dotcomfy installation? (y/n) ")
			fmt.Scan(&confirmation)
			LOGGER.Infof("Confirmation: %s", confirmation)

			if confirmation != "y" {
				LOGGER.Info("Aborting")
				fmt.Println("Aborting")
				os.Exit(0)
			}
		}

		user, err := user.Current()
		if err != nil {
			LOGGER.Fatal(err)
		}
		dotcomfy_dir := user.HomeDir + "/.dotcomfy"
		// Defaults to XDG_CONFIG_HOME if not set
		old_dotfiles_dir := user.HomeDir + "/.config"

		err = services.RemoveInstallation(dotcomfy_dir, old_dotfiles_dir)
		if err != nil {
			LOGGER.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uninstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uninstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	uninstallCmd.PersistentFlags().BoolVarP(&CONFIRM, "yes", "y", false, "Skips confirmation for uninstall")

}
