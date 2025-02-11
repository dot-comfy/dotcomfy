/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cobra

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dotcomfy",
	Short: "A simple tool for managing your dotfiles",
	Long: `A simple tool for managing your dotfiles.
	Whether you're SSHing into brand new cloud servers,
	bouncing between different operating systems, or 
	just wanting to try out different Linux rices,
	dotcomfy has you covered!`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&CFG_FILE, "config", "", "config file (default is $HOME/.dotcomfy/config.toml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if CFG_FILE != "" {
		viper.SetConfigFile(CFG_FILE)
	} else {
		cfg, err := os.UserConfigDir()
		cobra.CheckErr(err)

		CFG_FILE = cfg + "/dotcomfy/config.toml"

		fmt.Println("Using config file:", CFG_FILE)
	}
}
