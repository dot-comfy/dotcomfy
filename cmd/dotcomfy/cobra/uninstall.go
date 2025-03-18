/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cobra

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	Log "dotcomfy/internal/logger"
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

		// Delete symlinks and rename ".pre-dotcomfy" files back to their old names
		err = filepath.WalkDir(dotcomfy_dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				LOGGER.Error(err)
				return err
			}
			if !d.IsDir() {
				if !strings.Contains(path, ".git") && !strings.Contains(path, dotcomfy_dir+"README.md") {
					center_path := strings.TrimPrefix(path, dotcomfy_dir)
					old_path := old_dotfiles_dir + center_path
					if strings.Contains(old_path, ".pre-dotcomfy") {
						old_name := strings.Replace(old_path, ".pre-dotcomfy", "", 1)
						// Remove symlink
						err = os.Remove(old_path)
						if err != nil {
							LOGGER.Error(err)
							return err
						}
						err = os.Rename(old_path, old_name)
						if err != nil {
							LOGGER.Error(err)
							return err
						}
					} else { // Just remove symlink
						err = os.Remove(old_path)
						if err != nil {
							LOGGER.Error(err)
							return err
						}
					}
				}
			}
			return nil
		})

		preserve_config := true
		config_path := filepath.Join(dotcomfy_dir, "config.toml")
		_, err = os.Stat(config_path)
		if err == nil {
			preserve_config = false
		}

		// Delete everything in ~/.dotcomfy
		dir, err := os.Open(dotcomfy_dir)
		if err != nil {
			LOGGER.Error(err)
			os.Exit(1)
		}
		defer dir.Close()

		names, err := dir.Readdirnames(-1)
		if err != nil {
			LOGGER.Error(err)
			os.Exit(1)
		}

		for _, name := range names {
			file_path := filepath.Join(dotcomfy_dir, name)
			if preserve_config && file_path == config_path {
				continue
			}

			file_info, err := os.Stat(file_path)
			if err != nil {
				LOGGER.Error(err)
				os.Exit(1)
			}

			if file_info.IsDir() {
				err = os.RemoveAll(file_path)
				if err != nil {
					LOGGER.Error(err)
					continue
				}
			} else {
				err = os.Remove(file_path)
				if err != nil {
					LOGGER.Error(err)
					continue
				}
			}
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
