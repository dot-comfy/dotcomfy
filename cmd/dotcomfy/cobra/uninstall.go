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
)

var confirm string

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Removes current dotcomfy installation",
	Long: `Removes current dotcomfy installation
	from your system and restores previous dotfiles.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("uninstall called")
		var confirmation string
		fmt.Print("Are you sure you want to uninstall the current dotcomfy installation? (y/n)")
		fmt.Scan(&confirmation)

		if confirmation != "y" {
			fmt.Println("Aborting")
			os.Exit(0)
		}

		user, err := user.Current()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dotcomfy_dir := user.HomeDir + ".dotcomfy"
		old_dotfiles_dir := user.HomeDir + ".config"

		if len(args) > 0 {
			fmt.Println("Too many arguments")
			os.Exit(1)
		}

		// Delete symlinks and rename ".pre-dotcomfy" files back to their old names
		err = filepath.WalkDir(dotcomfy_dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "DEBUGPRINT: uninstall.go:36: err=%+v\n", err)
				return err
			}
			if !d.IsDir() {
				center_path := strings.TrimPrefix(path, dotcomfy_dir)
				old_path := old_dotfiles_dir + center_path
				if strings.Contains(old_path, ".pre-dotcomfy") {
					old_name := strings.Replace(old_path, ".pre-dotcomfy", "", 1)
					// Remove symlink
					err = os.Remove(old_path)
					if err != nil {
						fmt.Fprintf(os.Stderr, "DEBUGPRINT: uninstall.go:47: err=%+v\n", err)
						return err
					}
					err = os.Rename(old_path, old_name)
					if err != nil {
						fmt.Fprintf(os.Stderr, "DEBUGPRINT: uninstall.go:53: err=%+v\n", err)
						return err
					}
				} else { // Just remove symlink
					err = os.Remove(old_path)
					if err != nil {
						fmt.Fprintf(os.Stderr, "DEBUGPRINT: uninstall.go:60: err=%+v\n", err)
						return err
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
			fmt.Fprintf(os.Stderr, "DEBUGPRINT: uninstall.go:71: err=%+v\n", err)
			os.Exit(1)
		}
		defer dir.Close()

		names, err := dir.Readdirnames(-1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "DEBUGPRINT: uninstall.go:78: err=%+v\n", err)
			os.Exit(1)
		}

		for _, name := range names {
			file_path := filepath.Join(dotcomfy_dir, name)
			if preserve_config && file_path == config_path {
				continue
			}

			file_info, err := os.Stat(file_path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "DEBUGPRINT: uninstall.go:91: err=%+v\n", err)
				os.Exit(1)
			}

			if file_info.IsDir() {
				err = os.RemoveAll(file_path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "DEBUGPRINT: uninstall.go:98: err=%+v\n", err)
					continue
				}
			} else {
				err = os.Remove(file_path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "DEBUGPRINT: uninstall.go:103: err=%+v\n", err)
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
	rootCmd.PersistentFlags().BoolSliceVarP(&confirm, "yes", "y", false, "Skips confirmation for uninstall")

}
