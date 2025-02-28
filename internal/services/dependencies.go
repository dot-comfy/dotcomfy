package services

import (
	"errors"
	"fmt"
	// "os"
	// "os/exec"

	Config "dotcomfy/internal/config"
)

var errs []error

func InstallDependenciesLinux() error {
	Config.SetConfig()
	config := Config.GetConfig()
	errs := config.Validate()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		return errors.New("Invalid config file")
	}

	package_manager, err := CheckPackageManager()
	if err != nil {
		fmt.Println("Error getting package manager:", err)
		return err
	}

	// fmt.Println("Please enter your password to install dependencies...")
	// cmd := exec.Command("sudo", "-S", os.Args[0])
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err = cmd.Run()
	// if err != nil {
	// 	fmt.Println("Error running with sudo:", err)
	// 	return err
	// }

	for dependency := range config.Dependencies {
		d, err := Config.GetDependency(dependency)
		if err != nil {
			fmt.Println(err)
		}
		e := InstallDependency(d, package_manager)
		if e != nil {
			errs = append(errs, e...)
		}
	}
	return nil
}
