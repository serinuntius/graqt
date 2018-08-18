package viewer

import (
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Aggregates []string `yaml:"aggregates"`
}

func LoadConfig(reader io.Reader) (*Config, error) {
	dec := yaml.NewDecoder(reader)

	c := Config{}
	if err := dec.Decode(&c); err != nil {
		return nil, errors.Wrap(err, "Failed to Decode config")
	}

	return &c, nil
}
