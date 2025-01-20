/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"dotcomfy/cmd/dotcomfy/cobra"
	"fmt"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DEBUGPRINT[19]: main.go:13: err=%+v\n", err)
		os.Exit(1)
	}

	if user.Uid == "0" {
		fmt.Println("Running as sudo is not permitted, aborting")
		os.Exit(1)
	} else {
		cobra.Execute()
	}
}
