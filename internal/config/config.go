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
	PostInstallSteps  []string `toml:"post_install_steps,omitempty"`
	PostInstallScript string   `toml:"post_install_script,omitempty"`
	Steps             []string `toml:"steps,omitempty"`
	Script            string   `toml:"script,omitempty"`
	Version           string   `toml:"version,omitempty"`
}

var config Config

// TODO: viper is omitting dependencies with empty maps in the config file.
//
//	I'll probably have to parse the config file manually
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

// TODO: Validate dependencies in config.
//
//	If version exists, steps and/or script should not exist.
//	post_install_steps and/or post_install_script are mutually exclusive.
//	steps and script are mutually exclusive.
func ValidateConfig() error {

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
