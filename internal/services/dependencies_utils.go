package services

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"

	Config "dotcomfy/internal/config"
	Log "dotcomfy/internal/logger"
)

var (
	done      = false
	done_chan = make(chan struct{})
)

/*
 * This is super janky right now because there aren't any go libraries that
 * abstract away package management for every major package manager. I plan
 * to get this just working in the meantime, then go back and contribute to
 * @REF [syspkg](https://github.com/bluet/syspkg) to fill out the missing
 * package managers.
 */

func CheckPackageManager() (string, error) {
	exists := func(pm string) bool {
		_, err := exec.LookPath(pm)
		return err == nil
	}

	if exists("apt") {
		return "apt", nil
	} else if exists("dnf") {
		return "dnf", nil
	} else if exists("yum") {
		return "yum", nil
	} else if exists("yay") {
		return "yay", nil
	} else if exists("pacman") {
		return "pacman", nil
	} else if exists("zypper") {
		return "zypper", nil
	} else {
		return "", errors.New("Unknown package manager")
	}
}

func InstallDependency(d *Config.Dependency, pm string) []error {
	LOGGER = Log.GetLogger()
	var needs []string
	var errs []error

	needs = d.Needs
	if needs != nil {
		for _, need := range needs {
			LOGGER.Info("Need dependency \"" + need + "\" to install \"" + d.Name + "\"...")
			n, error := Config.GetDependency(need)
			if error != nil {
				fmt.Println(error)
				LOGGER.Error(error)
				err := errors.New("Error getting dependency \"" + need + "\"...")
				fmt.Println(err)
				LOGGER.Error(err)
				errs = append(errs, err)
				return errs
			}
			if n.FailedInstall {
				err := errors.New("Dependency \"" + need + "\" previously failed to install, skipping \"" + d.Name + "\"...")
				fmt.Println(err)
				LOGGER.Error(err)
				errs = append(errs, err)
				return errs
			}
			err := InstallDependency(n, pm)
			if err != nil {
				errs = append(errs, err...)
			}
		}
	}

	if d.Installed {
		LOGGER.Info("Dependency \"" + d.Name + "\" already installed, skipping...")
		return errs
	} else if d.GetFailedInstall() {
		err := errors.New("Dependency \"" + d.Name + "\" previously failed to install, skipping...")
		fmt.Println(err)
		LOGGER.Error(err)
		errs = append(errs, err)
		return errs
	} else if d.Version != "" {
		if d.Version == "latest" {
			err := InstallPackage(pm, d.Name, "")
			if err != nil {
				d.SetFailedInstall()
				fmt.Println("Dependency \"" + d.Name + "\" failed to install from package manager...")
				errs = append(errs, err)
			}
		} else {
			err := InstallPackage(pm, d.Name, d.Version)
			if err != nil {
				d.SetFailedInstall()
				fmt.Println("Dependency \"" + d.Name + "\" failed to install from package manager...")
				errs = append(errs, err)
			}
		}
		if d.PostInstallSteps != nil {
			err := HandleSteps(d.PostInstallSteps)
			if err != nil {
				d.SetFailedInstall()
				fmt.Println("Dependency \"" + d.Name + "\" failed during the post install steps...")
				errs = append(errs, err)
				return errs
			}
		} else if d.PostInstallScript != "" {
			// TODO: Handle post install script
		}
		d.SetInstalled()
	} else {
		fmt.Println("Installing dependency \"" + d.Name + "\"...")
		if d.Steps != nil {
			err := HandleSteps(d.Steps)
			if err != nil {
				d.SetFailedInstall()
				fmt.Println("Dependency \"" + d.Name + "\" failed during the install steps...")
				errs = append(errs, err)
				return errs
			}
		} else {
			// TODO: Handle script
		}
		d.SetInstalled()
	}
	return errs
}

func InstallPackage(pm string, pkg string, version string) error {
	LOGGER = Log.GetLogger()
	fmt.Println("Installing package \"" + pkg + "\" from package manager " + pm + " ...")
	LOGGER.Info("Installing package \"" + pkg + "\" from package manager " + pm + " ...")
	Progress()
	switch pm {
	case "apt":
		if version != "" {
			pkg = pkg + "=" + version
		}
		err := exec.Command("sudo", "apt", "install", "-y", pkg).Run()
		done = true
		return err
	case "dnf":
		if version != "" {
			pkg = pkg + "-" + version
		}
		cmd := fmt.Sprintf("sudo -S dnf install %s -y --skip-unavailable", pkg)
		command := exec.Command("/bin/sh", "-c", cmd)
		_, err := command.CombinedOutput()
		// fmt.Println(string(output))
		// LOGGER.Info(string(output))
		done = true
		return err
	case "yum":
		if version != "" {
			pkg = pkg + "=" + version
		}
		err := exec.Command("sudo", "yum", "install", "-y", pkg).Run()
		done = true
		return err
	case "pacman":
		if version != "" {
			pkg = pkg + "=" + version
		}
		cmd := fmt.Sprintf("sudo -S pacman -S %s --noconfirm", pkg)
		command := exec.Command("/bin/sh", "-c", cmd)
		_, err := command.CombinedOutput()
		done = true
		return err
	case "yay":
		if version != "" {
			log.Output(1, "Version not supported for yay")
		}
		err := exec.Command("yay", "--noconfirm", pkg).Run()
		done = true
		return err
	case "zypper":
		if version != "" {
			pkg = pkg + "=" + version
		}
		err := exec.Command("sudo", "zypper", "install", "-y", pkg).Run()
		done = true
		return err
	default:
		done = true
		return errors.New("Unknown package manager")
	}
}

func HandleSteps(steps []string) error {
	Progress()
	for _, step := range steps {
		cmd := exec.Command("/bin/sh", "-c", step)
		_, err := cmd.CombinedOutput()
		// fmt.Println(string(output))
		// LOGGER.Info(string(output))
		if err != nil {
			return err
		}
	}
	done = true
	return nil
}

func Progress() {
	fmt.Println("Entered progress")
	go func() {
		for {
			select {
			case <-time.After(100 * time.Millisecond):
				fmt.Print(".")
				if done == true {
					close(done_chan)
					return
				}
			case <-done_chan:
				done = false
				fmt.Printf(" Done\n")
				return
			}
		}
	}()
}
