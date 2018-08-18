package viewer

import (
	"regexp"

	"github.com/pkg/errors"
)

type Setting struct {
	AggregateRegexps []*regexp.Regexp
}

func newSetting() *Setting {
	return &Setting{}
}

func LoadSetting(c *Config) (*Setting, error) {
	s := newSetting()

	for _, ag := range c.Aggregates {
		re, err := regexp.Compile(ag)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to regexp.Compile. Aggregate: %s", ag)
		}
		s.AggregateRegexps = append(s.AggregateRegexps, re)
	}

	return s, nil
}
