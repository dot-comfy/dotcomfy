/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cobra

import (
	// "errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"dotcomfy/internal/services"
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

func run(cmd *cobra.Command, args []string) {
	fmt.Println("install called")
	user, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dotcomfy_dir := user.HomeDir + "/.dotcomfy"
	// Default to XDG_CONFIG_HOME directory if not set
	old_dotfiles_dir := user.HomeDir + "/.config"

	if len(args) > 1 {
		fmt.Println("Too many arguments")
		os.Exit(1)
	}

	os.MkdirAll(dotcomfy_dir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DEBUGPRINT: install.go:47: err=%+v\n", err)
		os.Exit(1)
	}

	if strings.Contains(args[0], "https://") {
		fmt.Println("Custom repo")
		var url string
		if !strings.HasSuffix(args[0], ".git") {
			url = args[0] + ".git"
		} else {
			url = args[0]
		}
		err = services.Clone(url, branch, dotcomfy_dir)

		if err != nil {
			fmt.Fprintf(os.Stderr, "DEBUGPRINT: install.go:58: err=%+v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Username")
		url := fmt.Sprintf("https://github.com/%s/dotfiles.git", args[0])
		err = services.Clone(url, branch, dotcomfy_dir)
		fmt.Fprintf(os.Stderr, "DEBUGPRINT: install.go:65: err=%+v\n", err)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// Walk through the cloned repo and perform rename/symlink operations
	err = filepath.WalkDir(dotcomfy_dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "DEBUGPRINT: install.go:78: err=%+v\n", err)
			return err
		}

		if !d.IsDir() {
			if strings.Contains(path, ".git") {
				fmt.Println("Skipping .git directory")
			} else if strings.Contains(path, dotcomfy_dir+"README.md") {
				fmt.Println("Skipping root level README.md")
			} else {
				_, err = services.RenameSymlinkUnix(old_dotfiles_dir, dotcomfy_dir, path)
			}
		}
		return nil
	})
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
	installCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "main", "Branch to clone")
}
