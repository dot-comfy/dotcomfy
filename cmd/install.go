/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	// "errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("install called")
		user, err := user.Current()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dotcomfy_path := user.HomeDir + ".dotcomfy"
		home_dir := user.HomeDir

		if len(args) > 1 {
			fmt.Println("Too many arguments")
			os.Exit(1)
		}

		temp_dir, err := os.MkdirTemp("", "dotcomfy-clone")
		if err != nil {
			fmt.Fprintf(os.Stderr, "DEBUGPRINT[3]: install.go:63: err=%+v\n", err)
			os.Exit(1)
		}

		defer os.RemoveAll(temp_dir)

		if strings.Contains(args[0], "dotfiles.git") {
			fmt.Println("Custom repo")
			_, err = git.PlainClone(temp_dir, false, &git.CloneOptions{
				URL: args[0],
			})

			if err != nil {
				fmt.Fprintf(os.Stderr, "DEBUGPRINT[2]: install.go:56: err=%+v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Username")
			url := fmt.Sprintf("https://github.com/%s/dotfiles.git", args[0])
			_, err = git.PlainClone(temp_dir, false, &git.CloneOptions{
				URL: url,
			})
			fmt.Fprintf(os.Stderr, "DEBUGPRINT[4]: install.go:70: err=%+v\n", err)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// Walk through the cloned repo and perform rename/symlink operations
		err = filepath.WalkDir(dotcomfy_path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "DEBUGPRINT[5]: install.go:84: err=%+v\n", err)
				return err
			}

			if !d.IsDir() {
				if strings.Contains(path, ".git") {
					fmt.Println("Skipping .git directory")
				} else if strings.Contains(path, temp_dir+"README.md") {
					fmt.Println("Skipping root level README.md")
				} else {

				}
			}
			return nil
		})
	},
}

func rename_symlink_unix(old_dotfiles_path, dotcomfy_path, new_path string) string {
	// TODO: implement this
	return ""
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
