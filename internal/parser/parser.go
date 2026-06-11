package parser

import "gopkg.in/yaml.v3"

func ParseConfig(data []byte) (any, error) {
	var cfg any
	err := yaml.Unmarshal(data, &cfg)
	return cfg, err
}
