/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cobra

import (
	"fmt"
	"os"
	"os/user"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch [OPTION]... NAME...",
	Short: "Switch to a different branch or repository of dotfiles",
	Long: `Switch to a different branch or repository of dotfiles.
	This will uninstall your current dotcomfy installation,
	and install the dotfiles from the branch or repository
	you are switching to.`,
	Args: cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if (branch == "") && (repo == "") || (branch != "" && repo != "") {
			return fmt.Errorf("Exactly one of --branch or --repo must be specified")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("switch called")
		user, err := user.Current()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dotcomfy_dir := user.HomeDir + "/.dotcomfy"
		// Defaults to XDG_CONFIG_HOME if not set
		old_dotfiles_dir := user.HomeDir + "/.config"

		if branch != "" {
			err = switchBranch(dotcomfy_dir, old_dotfiles_dir, branch)
			if err != nil {
				fmt.Fprintf(os.Stderr, "DEBUGPRINT: switch.go:41: err=%+v\n", err)
				os.Exit(1)
			}
		} else {
			err = switchRepo(dotcomfy_dir, old_dotfiles_dir, repo)
			if err != nil {
				fmt.Fprintf(os.Stderr, "DEBUGPRINT: switch.go:48: err=%+v\n", err)
				os.Exit(1)
			}
		}
	},
}

func switchBranch(dotcomfy_dir, old_dotfiles_dir, branch string) error {
	var repo_url string
	r, err := git.PlainOpen(dotcomfy_dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DEBUGPRINT: switch.go:59: err=%+v\n", err)
		return err
	}

	remote, err := r.Remote("origin")
	if err != nil {
		fmt.Fprintf(os.Stderr, "DEBUGPRINT: switch.go:64: err=%+v\n", err)
		return err
	}

	urls := remote.Config().URLs
	if len(urls) > 0 {
		repo_url = urls[0]
	} else {
		fmt.Println("No URL found for the remote 'origin'")
	}
}

func switchRepo(dotcomfy_dir, old_dotfiles_dir, repo string) error {

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
	switchCmd.Flags().StringVar(&branch, "branch", "", "Branch to switch to")
	switchCmd.Flags().StringVar(&repo, "repo", "", "Repository to switch to")
}
