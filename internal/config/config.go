package config

type Config struct {
	Dependencies map[string]Dependency `toml:"dependencies,omitempty"`
}

type Dependency struct {
	PostInstallSteps []string `toml:"post_install_steps,omitempty"`
	Steps            []string `toml:"steps,omitempty"`
	Script           string   `toml:"script,omitempty"`
	Version          string   `toml:"version,omitempty"`
}

var config Config

func GetConfig() Config {
	return config
}

func SetConfig(newConfig Config) {
	config = newConfig
}
