package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Dependencies map[string]Dependency `toml:"dependencies,omitempty"`
}

// TODO: Add "needs" cyclical dependency check
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

		if version != "" && steps != nil {
			errs = append(errs, errors.New("Dependency \""+dependency+"\" cannot have both \"version\" and \"steps\""))
		} else if version != "" && script != "" {
			errs = append(errs, errors.New("Dependency \""+dependency+"\" cannot have both \"version\" and \"script\""))
		} else if steps != nil && script != "" {
			errs = append(errs, errors.New("Dependency \""+dependency+"\" cannot have both \"steps\" and \"script\""))
		} else if post_install_steps != nil && post_install_script != "" {
			errs = append(errs, errors.New("Dependency \""+dependency+"\" cannot have both \"post_install_steps\" and \"post_install_script\""))
		}
	}
	return errs
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
	d.Installed = true
}

func (d *Dependency) GetInstalled() bool {
	return d.Installed
}

func (d *Dependency) SetFailedInstall() {
	d.FailedInstall = true
}

func (d *Dependency) GetFailedInstall() bool {
	return d.FailedInstall
}

var config Config

func SetConfig() {
	cfg, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DEBUGPRINT: config.go:27: err=%+v\n", err)
		os.Exit(1)
	}
	viper.AddConfigPath(cfg + "/dotcomfy/") // Config file lives in $HOME/.config/dotcomfy/
	viper.SetConfigName("config.toml")
	viper.SetConfigType("toml")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DEBUGPRINT: config.go:29: err=%+v\n", err)
	}
	viper.Unmarshal(&config)
	config.SetDependencyNames()
}

func GetConfig() *Config {
	return &config
}

func GetDependencies() map[string]string {
	dependencies := viper.GetStringMapString("dependencies")
	return dependencies
}
