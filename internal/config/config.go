package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type BuildType struct {
	Name  string   `yaml:"name"`
	Files []string `yaml:"files"`
}

type Types struct {
	BuildTypes []BuildType `yaml:"types"`
}

func LoadConfig(filename string) (Types, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Types{}, err
	}

	var t Types
	err = yaml.Unmarshal(bytes, &t)
	if err != nil {
		return Types{}, err
	}

	return t, nil
}
