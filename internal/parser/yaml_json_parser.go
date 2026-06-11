package parser

import "gopkg.in/yaml.v3"

type YamlJsonParser struct{}

func (p YamlJsonParser) ParseConfig(data []byte) (any, error) {
	var cfg any
	// парсит как yaml так и json
	err := yaml.Unmarshal(data, &cfg)
	return cfg, err
}
