package config

import "flag"

type Config struct {
	IsSilent bool
	IsStdin  bool
}

func Init() *Config {
	cfg := &Config{}

	flag.BoolVar(&cfg.IsSilent, "s", false, "do not exit with an error if there are errors")
	flag.BoolVar(&cfg.IsStdin, "stdin", false, "read config from stdin")

	flag.Parse()

	return cfg
}
