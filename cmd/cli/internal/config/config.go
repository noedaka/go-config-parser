package config

import "flag"

type Config struct {
	IsSilent    bool
	IsStdin     bool
	IsRecursive bool
}

func Init() *Config {
	cfg := &Config{}

	flag.BoolVar(&cfg.IsSilent, "s", false, "do not exit with an error if there are any errors")
	flag.BoolVar(&cfg.IsStdin, "stdin", false, "read config from stdin")
	flag.BoolVar(&cfg.IsRecursive, "r", false, "recursive directory analysis")

	flag.Parse()

	return cfg
}
