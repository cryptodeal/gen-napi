package config

import (
	"log"
	"os"

	"github.com/cryptodeal/gen-napi/napi"
	"gopkg.in/yaml.v2"
)

func ReadFromFilepath(cfgFilepath string) napi.Config {
	b, err := os.ReadFile(cfgFilepath)
	if err != nil {
		log.Fatalf("Could not read config file from %s: %v", cfgFilepath, err)
	}
	conf := napi.Config{}
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		log.Fatalf("Could not parse config file froms: %v", err)
	}

	return conf
}
