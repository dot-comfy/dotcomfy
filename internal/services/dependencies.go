package services

import (
	"errors"

	Config "dotcomfy/internal/config"
	Log "dotcomfy/internal/logger"
)

var errs []error

func InstallDependenciesLinux(pm string) error {
	LOGGER = Log.GetLogger()
	Config.SetConfig()
	config := Config.GetConfig()
	errs := config.Validate()
	if len(errs) > 0 {
		for _, err := range errs {
			LOGGER.Error(err)
		}
		return errors.New("Invalid config file")
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
			LOGGER.Error(err)
		}
		e := InstallDependency(d, pm)
		if e != nil {
			errs = append(errs, e...)
		}
	}
	return nil
}
