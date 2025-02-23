package services

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	Config "dotcomfy/internal/config"
)

var dependencies map[string]Config.Dependency
var errs []error

func InstallDependenciesLinux(config Config.Config) error {
	dependencies = config.Dependencies

	fmt.Println(dependencies)

	errs := Config.ValidateDependencies(dependencies)
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		return errors.New("Invalid config file")
	}

	package_manager, err := CheckPackageManager()

	fmt.Println("Please enter your password to install dependencies...")
	cmd := exec.Command("sudo", "-S", os.Args[0])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running with sudo:", err)
		return err
	}

	for dependency := range dependencies {
		e := InstallDependency(dependency, package_manager)
		if e != nil {
			errs = append(errs, e...)
		}
		/*
			for _, err := range errs {
				fmt.Println(err)
			}
		*/
		/*
			dependency_map, err := Config.GetDependency(dependency)
			d := dependencies[dependency]
			fmt.Printf("%s needs: %v\n", dependency, d.Needs)

			if err != nil {
				return err
			}

			// Handle version case first before any post install steps or scripts
			for k, v := range dependency_map {
				if k == "version" {
					fmt.Printf("Installing %s at version %s from package manager %s...\n", dependency, v.(string), package_manager)
					if v.(string) == "latest" {
						err = InstallPackage(package_manager, dependency, "")
					} else {
						err = InstallPackage(package_manager, dependency, v.(string))
					}
					if err != nil {
						fmt.Println("Error installing package:", err)
					}
				}
			}

			for k, v := range dependency_map {
				switch k {
				case "steps":
					fmt.Printf("Running steps for %s...\n", dependency)
					HandleSteps(v.([]interface{}))
				case "post_install_steps":
					for i, step := range v.([]interface{}) {
						fmt.Printf("Post install Step %d: %s\n", i, step)
					}
					case "script":
						fmt.Printf("Script location: %s\n", v)
						file_path := filepath.Join(dotcomfy_dir, v.(string))
						_, err := os.Stat(file_path)
						if os.IsNotExist(err) {
							fmt.Printf("File %s does not exist. Ensure it's in the directory %s\n", v, dotcomfy_dir)
							continue
						}
						if err != nil {
							fmt.Fprintf(os.Stderr, "DEBUGPRINT: dependencies.go:41: err=%+v\n", err)
							continue
						}
						content, err := os.ReadFile(file_path)
						if err != nil {
							fmt.Fprintf(os.Stderr, "DEBUGPRINT: dependencies.go:46: err=%+v\n", err)
							continue
						}
						fmt.Printf("Script content: %s\n", string(content))
					case "post_install_script":
						fmt.Printf("Post install Script location: %s\n", v)
						file_path := filepath.Join(dotcomfy_dir, v.(string))
						_, err := os.Stat(file_path)
						if os.IsNotExist(err) {
							fmt.Printf("File %s does not exist. Ensure it's in the directory %s\n", v, dotcomfy_dir)
							continue
						}
						if err != nil {
							fmt.Fprintf(os.Stderr, "DEBUGPRINT: dependencies.go:59: err=%+v\n", err)
							continue
						}
						content, err := os.ReadFile(file_path)
						if err != nil {
							fmt.Fprintf(os.Stderr, "DEBUGPRINT: dependencies.go:69: err=%+v\n", err)
							continue
						}
						fmt.Printf("Post Install Script content: %s\n", string(content))
				default:
					fmt.Printf("Unknown key: %s\n", k)
				}
			}

			if len(dependency_map) == 0 {
				fmt.Println("Installing package at latest version from package manager")
			}
		*/
	}
	return nil
}
