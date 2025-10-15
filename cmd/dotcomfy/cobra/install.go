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

	Log "dotcomfy/internal/logger"
	"dotcomfy/internal/services"
)

var skip_dependencies bool

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
	LOGGER = Log.GetLogger()
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
		LOGGER.Error(err)
		os.Exit(1)
	}

	if strings.Contains(args[0], "https://") {
		var url string
		if !strings.HasSuffix(args[0], ".git") {
			url = args[0] + ".git"
		} else {
			url = args[0]
		}
		err = services.Clone(url, BRANCH, COMMIT, dotcomfy_dir)

		if err != nil {
			LOGGER.Error(err)
			os.Exit(1)
		}
	} else {
		url := fmt.Sprintf("https://github.com/%s/dotfiles.git", args[0])
		err = services.Clone(url, BRANCH, COMMIT, dotcomfy_dir)
		if err != nil {
			LOGGER.Error(err)
		}

		if err != nil {
			LOGGER.Fatal(err)
		}
	}

	// Walk through the cloned repo and perform rename/symlink operations
	err = filepath.WalkDir(dotcomfy_dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			LOGGER.Error(err)
			return err
		}

		if !d.IsDir() && !strings.Contains(path, ".git") && !strings.Contains(path, "README.md") {
			// center_path represents the path of the directory entry
			// with the dotcomfy_path prefix removed.
			center_path := strings.TrimPrefix(path, dotcomfy_dir)
			_, err = services.RenameSymlinkUnix(old_dotfiles_dir, dotcomfy_dir, center_path)
			if err != nil {
				LOGGER.Error(err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		LOGGER.Error(err)
	}

	if !skip_dependencies {
		err = services.InstallDependenciesLinux()

		if err != nil {
			LOGGER.Error(err)
		}
	}

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
	installCmd.PersistentFlags().StringVarP(&BRANCH, "branch", "b", "main", "Branch to clone")
	installCmd.PersistentFlags().StringVar(&COMMIT, "at-commit", "", "Specific commit hash to install")
	installCmd.Flags().BoolVar(&skip_dependencies, "skip-dependencies", false, "Skip installing dependencies")
}
