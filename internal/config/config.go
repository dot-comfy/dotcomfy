package config

type Config struct {
	Foo string
}

var config Config

func GetConfig() Config {
	return config
}

func SetConfig(newConfig Config) {
	config = newConfig
}
