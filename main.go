/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"dotcomfy/cmd/dotcomfy/cobra"
	"fmt"
	"os"
)

func main() {
	if os.Geteuid() == 0 && os.Getenv("SUDO_USER") == "" {
		fmt.Println("Running as sudo is not permitted, aborting")
		os.Exit(1)
	} else {
		cobra.Execute()
	}
}
