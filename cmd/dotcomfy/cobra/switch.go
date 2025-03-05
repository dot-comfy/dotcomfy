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

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"

	Log "dotcomfy/internal/logger"
	"dotcomfy/internal/services"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch [OPTION]... NAME...",
	Short: "Switch to a different branch or repository of dotfiles",
	Long: `Switch to a different branch or repository of dotfiles.
	This will uninstall your current dotcomfy installation,
	and install the dotfiles from the branch or repository
	you are switching to.`,
	Args: cobra.MinimumNArgs(0),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if (BRANCH == "") && (REPO == "") {
			return fmt.Errorf("At least one of --branch or --repo must be specified")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		LOGGER = Log.GetLogger()
		var repo_url string
		user, err := user.Current()
		if err != nil {
			LOGGER.Fatal(err)
		}
		dotcomfy_dir := user.HomeDir + "/.dotcomfy"
		// Defaults to XDG_CONFIG_HOME if not set
		old_dotfiles_dir := user.HomeDir + "/.config"

		// Changing to different branch of same repo
		if REPO == "" && BRANCH != "" {
			r, err := git.PlainOpen(dotcomfy_dir)
			if err != nil {
				LOGGER.Errorf("switch.go:44: err=%+v\n", err)
			}

			remote, err := r.Remote("origin")
			if err != nil {
				LOGGER.Errorf("switch.go:49: err=%+v\n", err)
			}

			urls := remote.Config().URLs
			if len(urls) > 0 {
				repo_url = urls[0]
			} else {
				LOGGER.Fatal("No URL found for the remote 'origin'")
			}

			err = switchDotfiles(dotcomfy_dir, old_dotfiles_dir, repo_url, BRANCH)
			if err != nil {
				LOGGER.Fatalf("switch.go:61: err=%+v\n", err)
			}
		} else { // Changing to different repo
			if strings.Contains(REPO, "https://") {
				if !strings.HasSuffix(REPO, ".git") {
					repo_url = REPO + ".git"
				} else {
					repo_url = REPO
				}
			} else {
				repo_url = fmt.Sprintf("https://github.com/%s/dotfiles.git", REPO)
			}
			LOGGER.Errorf("switch.go:79: repo_url=%+v\n", repo_url)
			err = switchDotfiles(dotcomfy_dir, old_dotfiles_dir, repo_url, BRANCH)
			if err != nil {
				LOGGER.Fatalf("switch.go:67: err=%+v\n", err)
			}
		}
	},
}

func switchDotfiles(dotcomfy_dir, old_dotfiles_dir, url, branch string) error {
	// Perform uninstall
	err := services.RemoveInstallation(dotcomfy_dir, old_dotfiles_dir)
	if err != nil {
		LOGGER.Errorf("switch.go:92: err=%+v\n", err)
		return err
	}

	// Perform install
	err = services.Clone(url, branch, dotcomfy_dir)
	if err != nil {
		LOGGER.Errorf("switch.go:99: err=%+v\n", err)
		return err
	}

	// Walk through the cloned repo and perform rename/symlink operations
	err = filepath.WalkDir(dotcomfy_dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			LOGGER.Errorf("switch.go:106: err=%+v\n", err)
			return err
		}

		if !d.IsDir() {
			if !strings.Contains(path, ".git") && !strings.Contains(path, dotcomfy_dir+"README.md") {
				_, err = services.RenameSymlinkUnix(old_dotfiles_dir, dotcomfy_dir, path)
				if err != nil {
					LOGGER.Errorf("switch.go:114: err=%+v\n", err)
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		LOGGER.Errorf("switch.go:122: err=%+v\n", err)
		return err
	}
	return nil
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
	switchCmd.Flags().StringVar(&BRANCH, "branch", "", "Branch to switch to")
	switchCmd.Flags().StringVar(&REPO, "repo", "", "Repository to switch to")
}
