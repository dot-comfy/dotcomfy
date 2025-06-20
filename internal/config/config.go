package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	Log "dotcomfy/internal/logger"
)

type Config struct {
	Dependencies map[string]Dependency `toml:"dependencies,omitempty"`
	Auth         Auth                  `toml:authentication,omitempty`
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

type Auth struct {
	Username    string `toml:username,omitempty`
	Email       string `toml:email,omitempty`
	SSHFilePath string `toml:ssh_file,omitempty`
}

func (g *Auth) GetUsername() string {
	return g.Username
}

func (g *Auth) GetEmail() string {
	return g.Email
}

func (g *Auth) GetSSHFilePath() string {
	return g.SSHFilePath
}

type Dependency struct {
	Name              string   `toml:"name,omitempty"`
	Needs             []string `toml:"needs,omitempty"`
	PostInstallSteps  []string `toml:"post_install_steps,omitempty"`
	PostInstallScript string   `toml:"post_install_script,omitempty"`
	Steps             []string `toml:"steps,omitempty"`
	Script            string   `toml:"script,omitempty"`
	Version           string   `toml:"version,omitempty"`
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

var config *Config

func SetConfig() {
	LOGGER := Log.GetLogger()
	cfg, err := os.UserConfigDir()
	if err != nil {
		LOGGER.Fatal(err)
	}
	viper.AddConfigPath(cfg + "/dotcomfy/") // Config file lives in $HOME/.config/dotcomfy/
	viper.SetConfigName("config.toml")
	viper.SetConfigType("toml")
	err = viper.ReadInConfig()
	if err != nil {
		LOGGER.Error(err)
	}
	viper.Unmarshal(&config)
	config.SetDependencyNames()
}

func GetConfig() *Config {
	return config
}

func GetDependencies() map[string]string {
	dependencies := viper.GetStringMapString("dependencies")
	return dependencies
}
