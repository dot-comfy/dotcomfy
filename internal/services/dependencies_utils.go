package services

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

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
	/*
		file, err := os.Open(XDG_CONFIG_HOME + "/dotcomfy/" + file_name)
		if err != nil {
			LOGGER.Error("Error opening the file \""+file_name+"\":", err)
			return err
		}
		defer file.Close()

		var lines []string
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
			if line != "" {
				lines = append(lines, line)
			}
		}

		if err := scanner.Err(); err != nil {
			LOGGER.Error("Error occurred during file scanning:", err)
		}

		for _, line := range lines {
			cmd := exec.Command("/bin/sh", "-c", line)
			output, err := cmd.CombinedOutput()
			// fmt.Println(string(output))
			LOGGER.Info(string(output))
			if err != nil {
				return err
			}
		}
	*/
	content, err := os.ReadFile(XDG_CONFIG_HOME + "/dotcomfy/" + file_name)
	if err != nil {
		LOGGER.Error("Error opening the file \""+file_name+"\":", err)
		return err
	}

	cmd := exec.Command("/bin/sh", "-c", string(content))

	std_output, err := cmd.StdoutPipe()
	if err != nil {
		LOGGER.Error("Error creating stdoutpipe:", err)
		return err
	}

	std_err, err := cmd.StderrPipe()
	if err != nil {
		LOGGER.Error("Error creating stderrpipe:", err)
		return err
	}

	std_out_bytes, _ := io.ReadAll(std_output)
	std_err_bytes, _ := io.ReadAll(std_err)

	if err := cmd.Wait(); err != nil {
		LOGGER.Errorf("\""+file_name+"\" finished with an error:", err)
		return err
	}

	_ = string(std_out_bytes)

	if len(std_err_bytes) > 0 {
		LOGGER.Errorf("Errors:", string(std_err_bytes))
	}

	return nil
}
