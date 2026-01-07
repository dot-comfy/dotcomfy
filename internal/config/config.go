package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	Log "dotcomfy/internal/logger"
)

type Config struct {
	Dependencies   map[string]Dependency `yaml:"dependencies,omitempty"`
	Auth           Auth                  `yaml:"authentication,omitempty"`
	PackageManager string                `yaml:"package_manager,omitempty"`
}

// TODO: Find a way to pull config file down first from the repo if it exists to validate before installation
func (c *Config) Validate() []error {
	dependencies := c.Dependencies
	errs := []error{}
	for dependency := range dependencies {
		d, err := GetDependency(dependency)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		version := d.GetVersion()
		post_install_steps := d.GetPostInstallSteps()
		post_install_script := d.GetPostInstallScript()
		steps := d.GetSteps()
		script := d.GetScript()

		// Check for conflicting fields
		if version != "" && steps != nil {
			errs = append(errs, errors.New("Dependency \""+dependency+"\" cannot have both \"version\" and \"steps\""))
		}
		if version != "" && script != "" {
			errs = append(errs, errors.New("Dependency \""+dependency+"\" cannot have both \"version\" and \"script\""))
		}
		if steps != nil && script != "" {
			errs = append(errs, errors.New("Dependency \""+dependency+"\" cannot have both \"steps\" and \"script\""))
		}
		if post_install_steps != nil && post_install_script != "" {
			errs = append(errs, errors.New("Dependency \""+dependency+"\" cannot have both \"post_install_steps\" and \"post_install_script\""))
		}
		if post_install_steps != nil && version == "" {
			errs = append(errs, errors.New("Dependency \""+dependency+"\" must have \"version\" if using \"post_install_steps\""))
		}
		if post_install_script != "" && version == "" {
			errs = append(errs, errors.New("Dependency \""+dependency+"\" must have \"version\" if using \"post_install_script\""))
		}
		if version == "" && script == "" && steps == nil && post_install_steps == nil && post_install_script == "" {
			errs = append(errs, errors.New("Dependency \""+dependency+"\" must have \"version\" set to \"latest\" or a specific version number"))
		}

		if d.GetNeeds() != nil {
			for _, n := range d.GetNeeds() {
				if n == dependency {
					errs = append(errs, errors.New("Dependency \""+dependency+"\" cannot have itself as a \"need\""))
				} else {
					// Check to see if there is a "needs" cycle
					fmt.Println("Checking dependency \"" + dependency + "\" for a dependency cycle...")
					cycle, chain := CheckDependencyCycle(dependency, n)
					chain = append(chain, n)
					fmt.Println(chain)
					if cycle {
						errs = append(errs, errors.New("Dependency \""+dependency+"\" has a dependency cycle: "+strings.Join(chain, " <- ")+" <- "+dependency))
					}
				}
			}
		}
	}
	if c.PackageManager != "" && !isValidPackageManager(c.PackageManager) {
		errs = append(errs, errors.New("Invalid preferred_package_manager: "+c.PackageManager+". Must be one of: apt, dnf, yum, yay, pacman, zypper, brew"))
	}
	return errs
}

func CheckDependencyCycle(dependency string, need string) (bool, []string) {
	d, err := GetDependency(need)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	if d.GetNeeds() != nil {
		for _, n := range d.GetNeeds() {
			fmt.Println(n)
			if n == dependency {
				return true, []string{n}
			} else {
				cycle, chain := CheckDependencyCycle(dependency, n)
				if cycle {
					return true, append(chain, n)
				}
			}
		}
	}
	return false, nil
}

func GetDependency(name string) (*Dependency, error) {
	c := GetConfig()
	for _, d := range c.Dependencies {
		if d.Name == name {
			return &d, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Dependency \"%s\" not found", name))
}

func (c *Config) SetDependencyNames() {
	newDependencies := make(map[string]Dependency)
	for name, dependency := range c.Dependencies {
		d := dependency
		d.Name = name
		newDependencies[name] = d
	}
	c.Dependencies = newDependencies
}

// TODO: I need to do error handling on these getters since they may not exist
type Auth struct {
	Username         string `yaml:"username,omitempty"`
	Email            string `yaml:"email,omitempty"`
	SSHFile          string `yaml:"ssh_file,omitempty"`
	SSHKeyPassphrase string `yaml:"ssh_key_passphrase,omitempty"`
}

func (g *Auth) GetUsername() string {
	return g.Username
}

func (g *Auth) GetEmail() string {
	return g.Email
}

func (g *Auth) GetSSHKeyPath() (string, error) {
	if strings.HasPrefix(g.SSHFile, "~/") || g.SSHFile == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if g.SSHFile == "~" {
			return home, nil
		}
		return filepath.Join(home, strings.TrimPrefix(g.SSHFile, "~/")), nil
	}
	return g.SSHFile, nil
}

func (g *Auth) GetSSHKeyPassphrase() string {
	return g.SSHKeyPassphrase
}

type Dependency struct {
	Name              string   `yaml:"name,omitempty"`
	Needs             []string `yaml:"needs,omitempty"`
	PostInstallSteps  []string `yaml:"post_install_steps,omitempty"`
	PostInstallScript string   `yaml:"post_install_script,omitempty"`
	Steps             []string `yaml:"steps,omitempty"`
	Script            string   `yaml:"script,omitempty"`
	Version           string   `yaml:"version,omitempty"`
	Installed         bool
	FailedInstall     bool
}

func (d *Dependency) SetName(name string) {
	d.Name = name
}

func (d *Dependency) GetName() string {
	return d.Name
}

func (d *Dependency) GetNeeds() []string {
	return d.Needs
}

func (d *Dependency) GetPostInstallSteps() []string {
	return d.PostInstallSteps
}

func (d *Dependency) GetPostInstallScript() string {
	return d.PostInstallScript
}

func (d *Dependency) GetSteps() []string {
	return d.Steps
}

func (d *Dependency) GetScript() string {
	return d.Script
}

func (d *Dependency) GetVersion() string {
	return d.Version
}

func (d *Dependency) SetInstalled() {
	Log.GetLogger().Info("Setting dependency \"" + d.Name + "\" as installed...")
	c := GetConfig()
	if dependency, ok := c.Dependencies[d.Name]; ok {
		dependency.Installed = true
		c.Dependencies[d.Name] = dependency
	}
}

func (d *Dependency) GetInstalled() bool {
	return d.Installed
}

func (d *Dependency) SetFailedInstall() {
	Log.GetLogger().Info("Setting dependency \"" + d.Name + "\" as failed install...")
	c := GetConfig()
	if dependency, ok := c.Dependencies[d.Name]; ok {
		dependency.FailedInstall = true
		c.Dependencies[d.Name] = dependency
	}
}

func (d *Dependency) GetFailedInstall() bool {
	return d.FailedInstall
}

func isValidPackageManager(pm string) bool {
	validPMs := []string{"apt", "brew", "dnf", "yum", "yay", "pacman", "zypper"}
	for _, v := range validPMs {
		if pm == v {
			return true
		}
	}
	return false
}

var config *Config

func SetConfig() {
	LOGGER := Log.GetLogger()
	cfg, err := os.UserHomeDir()

	// Create a new Viper instance to avoid any global state issues
	v := viper.New()
	v.AddConfigPath(cfg + "/.config/dotcomfy/") // Config file lives in $HOME/.config/dotcomfy/
	v.SetConfigName("config.yaml")
	v.SetConfigType("yaml")
	err = v.ReadInConfig()
	if err != nil {
		LOGGER.Error(err)
	}

	localConfig := &Config{}

	// Debug what this fresh Viper instance read
	LOGGER.Info("Fresh Viper all settings:", v.AllSettings())
	LOGGER.Info("Fresh Viper auth:", v.Get("authentication"))

	// Unmarshal dependencies first
	var dependencies map[string]Dependency
	if deps := v.Get("dependencies"); deps != nil {
		if depMap, ok := deps.(map[string]any); ok {
			dependencies = make(map[string]Dependency)
			for key, value := range depMap {
				if depStruct, ok := value.(map[string]any); ok {
					dep := Dependency{}
					for fieldKey, fieldValue := range depStruct {
						switch fieldKey {
						case "version":
							if version, ok := fieldValue.(string); ok {
								dep.Version = version
							}
						case "script":
							if script, ok := fieldValue.(string); ok {
								dep.Script = script
							}
						case "steps":
							if steps, ok := fieldValue.([]any); ok {
								for _, step := range steps {
									if stepStr, ok := step.(string); ok {
										dep.Steps = append(dep.Steps, stepStr)
									}
								}
							}
						case "post_install_script":
							if script, ok := fieldValue.(string); ok {
								dep.PostInstallScript = script
							}
						case "post_install_steps":
							if steps, ok := fieldValue.([]any); ok {
								for _, step := range steps {
									if stepStr, ok := step.(string); ok {
										dep.PostInstallSteps = append(dep.PostInstallSteps, stepStr)
									}
								}
							}
						case "needs":
							if needs, ok := fieldValue.([]any); ok {
								for _, need := range needs {
									if needStr, ok := need.(string); ok {
										dep.Needs = append(dep.Needs, needStr)
									}
								}
							}
						}
					}
					dependencies[key] = dep
				}
			}
		}
	}

	// Unmarshal authentication separately
	var auth Auth
	if authData := v.Get("authentication"); authData != nil {
		if authMap, ok := authData.(map[string]any); ok {
			for key, value := range authMap {
				switch key {
				case "username":
					if username, ok := value.(string); ok {
						auth.Username = username
					}
				case "email":
					if email, ok := value.(string); ok {
						auth.Email = email
					}
				case "ssh_file":
					if sshFile, ok := value.(string); ok {
						auth.SSHFile = sshFile
					}
				case "ssh_key_passphrase":
					if passphrase, ok := value.(string); ok {
						auth.SSHKeyPassphrase = passphrase
					}
				}
			}
		}
	}
	// Unmarshal preferred package manager
	var preferredPM string
	if pm := v.Get("preferred_package_manager"); pm != nil {
		if pmStr, ok := pm.(string); ok {
			preferredPM = pmStr
		}
	}
	localConfig.PackageManager = preferredPM

	// Combine into final config
	localConfig.Dependencies = dependencies
	localConfig.Auth = auth

	LOGGER.Info("Config after manual unmarshal:", localConfig)

	// Update global config
	config = localConfig

	config.SetDependencyNames()
}

func SetTempConfig(p string) {
	LOGGER := Log.GetLogger()

	// Create a new Viper instance to avoid any global state issues
	v := viper.New()
	v.AddConfigPath(p)
	v.SetConfigName("config.yaml")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		LOGGER.Error(err)
	}

	localConfig := &Config{}

	// Debug what this fresh Viper instance read
	LOGGER.Info("Fresh Viper all settings:", v.AllSettings())
	LOGGER.Info("Fresh Viper auth:", v.Get("authentication"))

	// Unmarshal dependencies first
	var dependencies map[string]Dependency
	if deps := v.Get("dependencies"); deps != nil {
		if depMap, ok := deps.(map[string]any); ok {
			dependencies = make(map[string]Dependency)
			for key, value := range depMap {
				if depStruct, ok := value.(map[string]any); ok {
					dep := Dependency{}
					for fieldKey, fieldValue := range depStruct {
						switch fieldKey {
						case "version":
							if version, ok := fieldValue.(string); ok {
								dep.Version = version
							}
						case "script":
							if script, ok := fieldValue.(string); ok {
								dep.Script = script
							}
						case "steps":
							if steps, ok := fieldValue.([]any); ok {
								for _, step := range steps {
									if stepStr, ok := step.(string); ok {
										dep.Steps = append(dep.Steps, stepStr)
									}
								}
							}
						case "post_install_script":
							if script, ok := fieldValue.(string); ok {
								dep.PostInstallScript = script
							}
						case "post_install_steps":
							if steps, ok := fieldValue.([]any); ok {
								for _, step := range steps {
									if stepStr, ok := step.(string); ok {
										dep.PostInstallSteps = append(dep.PostInstallSteps, stepStr)
									}
								}
							}
						case "needs":
							if needs, ok := fieldValue.([]any); ok {
								for _, need := range needs {
									if needStr, ok := need.(string); ok {
										dep.Needs = append(dep.Needs, needStr)
									}
								}
							}
						}
					}
					dependencies[key] = dep
				}
			}
		}
	}

	// Unmarshal authentication separately
	var auth Auth
	if authData := v.Get("authentication"); authData != nil {
		if authMap, ok := authData.(map[string]any); ok {
			for key, value := range authMap {
				switch key {
				case "ssh_file":
					if sshFile, ok := value.(string); ok {
						auth.SSHFile = sshFile
					}
					// Unmarshal preferred package manager
					var preferredPM string
					if pm := v.Get("preferred_package_manager"); pm != nil {
						if pmStr, ok := pm.(string); ok {
							preferredPM = pmStr
						}
					}
					localConfig.PackageManager = preferredPM
				case "ssh_key_passphrase":
					if passphrase, ok := value.(string); ok {
						auth.SSHKeyPassphrase = passphrase
					}
				}
			}
		}
	}

	// Combine into final config
	localConfig.Dependencies = dependencies
	localConfig.Auth = auth

	LOGGER.Info("Config after manual unmarshal:", localConfig)

	// Update global config
	config = localConfig

	config.SetDependencyNames()
}

func GetConfig() *Config {
	return config
}

func GetDependencies() map[string]string {
	dependencies := viper.GetStringMapString("dependencies")
	return dependencies
}
