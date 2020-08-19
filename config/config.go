package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config : application config stored as global variable
var Config *EnvironmentConfig

// EnvironmentConfig :
type EnvironmentConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

// DatabaseConfig :
type DatabaseConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Name string `yaml:"name"`
}

// ServerConfig :
type ServerConfig struct {
	Host   string       `yaml:"host"`
	Port   string       `yaml:"port"`
	Static StaticConfig `yaml:"static"`
}

// StaticConfig :
type StaticConfig struct {
	ManifestPath string `yaml:"manifestpath"`
}

// Init :
func Init() {
	Config = readConfig()
}

func readConfig() *EnvironmentConfig {
	file := fmt.Sprintf("config/environments/%s.yml", os.Getenv("AGILELY_ENV"))
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
	defer f.Close()

	var cfg EnvironmentConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
	return &cfg
}
