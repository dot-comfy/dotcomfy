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

var config Config

func GetConfig() Config {
	cfg, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DEBUGPRINT: config.go:27: err=%+v\n", err)
		os.Exit(1)
	}
	viper.AddConfigPath(cfg + "/dotcomfy/") // Config file lives in $HOME/.config/dotcomfy/
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DEBUGPRINT: config.go:29: err=%+v\n", err)
	}
	viper.Unmarshal(&config)
	fmt.Println(config)
	return config
}

func SetConfig(newConfig Config) {
	config = newConfig
}

// TODO: Add "needs" cyclical dependency check
func ValidateDependencies(dependencies map[string]Dependency) []error {
	errs := []error{}
	for dependency := range dependencies {
		dependency_map, err := GetDependency(dependency)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		_, version_exists := dependency_map["version"]
		_, post_install_steps_exists := dependency_map["post_install_steps"]
		_, post_install_script_exists := dependency_map["post_install_script"]
		_, steps_exists := dependency_map["steps"]
		_, script_exists := dependency_map["script"]

		if version_exists && steps_exists {
			errs = append(errs, errors.New("Dependency \""+dependency+"\" cannot have both \"version\" and \"steps\""))
		} else if version_exists && script_exists {
			errs = append(errs, errors.New("Dependency \""+dependency+"\" cannot have both \"version\" and \"script\""))
		} else if steps_exists && script_exists {
			errs = append(errs, errors.New("Dependency \""+dependency+"\" cannot have both \"steps\" and \"script\""))
		} else if post_install_steps_exists && post_install_script_exists {
			errs = append(errs, errors.New("Dependency \""+dependency+"\" cannot have both \"post_install_steps\" and \"post_install_script\""))
		}
	}
	return errs
}

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
