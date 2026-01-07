package services

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	Config "dotcomfy/internal/config"
	Log "dotcomfy/internal/logger"
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

func isValidPackageManager(pm string) bool {
	validPMs := []string{"apt", "dnf", "yum", "yay", "pacman", "zypper", "brew"}
	for _, v := range validPMs {
		if pm == v {
			return true
		}
	}
	return false
}

func findAvailablePM() string {
	exists := func(pm string) bool {
		_, err := exec.LookPath(pm)
		return err == nil
	}

	if exists("apt") {
		return "apt"
	} else if exists("dnf") {
		return "dnf"
	} else if exists("yum") {
		return "yum"
	} else if exists("yay") {
		return "yay"
	} else if exists("pacman") {
		return "pacman"
	} else if exists("zypper") {
		return "zypper"
	} else if exists("brew") {
		return "brew"
	} else {
		return ""
	}
}

func promptContinue(question string) bool {
	fmt.Print(question)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
	return answer == "y" || answer == "yes"
}

func ValidateAndGetPackageManager(preferredPM string) (string, error) {
	if preferredPM != "" {
		if !isValidPackageManager(preferredPM) {
			// invalid, find fallback
			fallback := findAvailablePM()
			if fallback == "" {
				return "", errors.New("Preferred package manager '" + preferredPM + "' is invalid and no fallback found")
			}
			// prompt
			if !promptContinue("The package manager '" + preferredPM + "' is not supported. Found '" + fallback + "' on your system. Continue with " + fallback + "? (y/N): ") {
				return "", errors.New("User declined to continue with fallback package manager")
			}
			return fallback, nil
		} else {
			// valid, check if exists
			exists := func(pm string) bool {
				_, err := exec.LookPath(pm)
				return err == nil
			}
			if exists(preferredPM) {
				return preferredPM, nil
			} else {
				// not found, find fallback
				fallback := findAvailablePM()
				if fallback == "" {
					return "", errors.New("Preferred package manager '" + preferredPM + "' not found and no fallback available")
				}
				if !promptContinue("Preferred package manager '" + preferredPM + "' not found. Found '" + fallback + "' on your system. Continue with " + fallback + "? (y/N): ") {
					return "", errors.New("User declined to continue with fallback package manager")
				}
				return fallback, nil
			}
		}
	} else {
		// no preferred, find available
		pm := findAvailablePM()
		if pm == "" {
			return "", errors.New("No package manager found")
		}
		return pm, nil
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
			err := HandleScript(d.PostInstallScript)
			if err != nil {
				d.SetFailedInstall()
				fmt.Println("Dependency \"" + d.Name + "\" failed during the install steps...")
				LOGGER.Error("Dependency \"" + d.Name + "\" failed during the install steps...")
				errs = append(errs, err)
				return errs
			}
		}
		d.SetInstalled()
	} else {
		fmt.Println("Installing dependency \"" + d.Name + "\"...")
		if d.Steps != nil {
			err := HandleSteps(d.Steps)
			if err != nil {
				d.SetFailedInstall()
				fmt.Println("Dependency \"" + d.Name + "\" failed during the install steps...")
				LOGGER.Error("Dependency \"" + d.Name + "\" failed during the install steps...")
				errs = append(errs, err)
				return errs
			}
		} else {
			err := HandleScript(d.Script)
			if err != nil {
				d.SetFailedInstall()
				fmt.Println("Dependency \"" + d.Name + "\" failed during the install steps...")
				LOGGER.Error("Dependency \"" + d.Name + "\" failed during the install steps...")
				errs = append(errs, err)
				return errs
			}
		}
		d.SetInstalled()
	}
	return errs
}

func InstallPackage(pm string, pkg string, version string) error {
	LOGGER = Log.GetLogger()

	fmt.Println("Installing package \"" + pkg + "\" from package manager " + pm + " ...")
	LOGGER.Info("Installing package \"" + pkg + "\" from package manager " + pm + " ...")

	switch pm {
	case "apt":
		if version != "" {
			pkg = pkg + "=" + version
		}
		err := exec.Command("sudo", "apt", "install", "-y", pkg).Run()
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
		return err
	case "yum":
		if version != "" {
			pkg = pkg + "=" + version
		}
		err := exec.Command("sudo", "yum", "install", "-y", pkg).Run()
		return err
	case "pacman":
		if version != "" {
			pkg = pkg + "=" + version
		}
		cmd := fmt.Sprintf("sudo -S pacman -S %s --noconfirm", pkg)
		command := exec.Command("/bin/sh", "-c", cmd)
		_, err := command.CombinedOutput()
		return err
	case "yay":
		if version != "" {
			log.Output(1, "Version not supported for yay")
		}
		err := exec.Command("yay", "--noconfirm", pkg).Run()
		return err
	case "zypper":
		if version != "" {
			pkg = pkg + "=" + version
		}
		err := exec.Command("sudo", "zypper", "install", "-y", pkg).Run()
		return err
	case "brew":
		if version != "" && version != "latest" {
			pkg = pkg + "@" + version
		}
		err := exec.Command("brew", "install", pkg).Run()
		return err
	default:
		return errors.New("Unknown package manager")
	}
}

func HandleSteps(steps []string) error {
	LOGGER = Log.GetLogger()

	for _, step := range steps {
		cmd := exec.Command("/bin/sh", "-c", step)
		output, err := cmd.CombinedOutput()
		// fmt.Println(string(output))
		LOGGER.Info(string(output))
		if err != nil {
			LOGGER.Error(err)
			return err
		}
	}
	return nil
}

func HandleScript(file_name string) error {
	LOGGER = Log.GetLogger()
	XDG_CONFIG_HOME, _ := os.UserConfigDir()

	err := os.Chmod(XDG_CONFIG_HOME+"/dotcomfy/"+file_name, 0755)
	if err != nil {
		LOGGER.Error("Error making script \""+file_name+"\" executable:", err)
		return err
	}

	cmd := exec.Command(XDG_CONFIG_HOME + "/dotcomfy/" + file_name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		LOGGER.Error("Error executing script \""+file_name+"\":", err)
		return err
	}

	fmt.Println(string(output))

	return nil
}
