package config

import (
	"errors"
	// "os"

	"github.com/spf13/viper"
)

type Config struct {
	Dependencies map[string]Dependency `toml:"dependencies,omitempty"`
}

type Dependency struct {
	PostInstallSteps  []string `toml:"post_install_steps,omitempty"`
	PostInstallScript string   `toml:"post_install_script,omitempty"`
	Steps             []string `toml:"steps,omitempty"`
	Script            string   `toml:"script,omitempty"`
	Version           string   `toml:"version,omitempty"`
}

/*
var config Config

func GetConfig() Config {
	data, err := os.ReadFile()
	return config
}

func SetConfig(newConfig Config) {
	config = newConfig
}
*/

func GetDependencies() map[string]string {
	dependencies := viper.GetStringMapString("dependencies")
	return dependencies
}

func GetDependency(name string) (map[string]interface{}, error) {
	exists := viper.IsSet("dependencies." + name)
	if !exists {
		return nil, errors.New("Dependency not found")
	}
	dependency := viper.Get("dependencies." + name).(map[string]interface{})
	return dependency, nil
}
