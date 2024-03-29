package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Prometheus struct {
		Outfile string `yaml:"outfile"`
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
	} `yaml:"prometheus"`
	SSH struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		User    string `yaml:"user"`
		Keyfile string `yaml:"keyfile"`
		Command string `yaml:"command"`
	} `yaml:"ssh"`
	LogLevel       string        `yaml:"loglevel"`
	ScrapeInterval time.Duration `yaml:"scrape_interval"`
}

func ReadConfig(filename string) *AppConfig {

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read file=%s err=%v ", filename, err)
	}
	var cfg AppConfig
	err = cfg.parse(yamlFile)

	if err != nil {
		log.Fatalf("Failed to unmarshal: %v", err)
	}
	fmt.Printf("%+v\n", cfg)
	return &cfg
}

func (c *AppConfig) parse(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}

	// Optional add checks for required fields

	return nil
}
