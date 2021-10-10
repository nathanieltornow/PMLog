package benchmark

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

type BenchConfig struct {
	Ops       int           `yaml:"ops"`
	Threads   int           `yaml:"threads"`
	Runtime   time.Duration `yaml:"runtime"`
	Endpoints []string      `yaml:"endpoints"`
}

func GetBenchConfig(configFile string) (*BenchConfig, error) {
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	b := &BenchConfig{}
	err = yaml.Unmarshal(buf, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
