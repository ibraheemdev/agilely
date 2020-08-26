package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"gopkg.in/yaml.v3"
)

// ReadConfig : read configuration files
func ReadConfig() (*engine.Config, error) {
	file := fmt.Sprintf("config/environments/%s.yml", os.Getenv("AGILELY_ENV"))
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &engine.Config{}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
	return cfg, nil
}
