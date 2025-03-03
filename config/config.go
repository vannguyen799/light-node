package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// GRPC configuration
	GRPC struct {
		URL string `yaml:"url"`
	} `yaml:"grpc"`

	// Cosmos configuration
	Cosmos struct {
		ContractAddr string `yaml:"contract_addr"`
	} `yaml:"cosmos"`

	// ZK Prover configuration
	ZKProver struct {
		URL string `yaml:"url"`
	} `yaml:"zk_prover"`

	// API authentication
	API struct {
		AuthToken string `yaml:"auth_token"`
		Timeout   int    `yaml:"timeout_seconds"`
	} `yaml:"api"`
}

var ConfigFilePath = flag.String(
	"c",
	"config.yml",
	"Specify the config path, default: 'config.yml' (root dir)",
)

func GetConfig() Config {
	var cfg Config
	flag.Parse()

	readFile(&cfg)
	validateConfig(&cfg)

	return cfg
}

func validateConfig(cfg *Config) {
	if cfg.GRPC.URL == "" {
		log.Fatal("GRPC URL is required in config file")
	}

	if cfg.Cosmos.ContractAddr == "" {
		log.Fatal("Cosmos contract address is required in config file")
	}

	if cfg.ZKProver.URL == "" {
		log.Fatal("ZK Prover URL is required in config file")
	}

	if cfg.API.Timeout == 0 {
		cfg.API.Timeout = 100 // default timeout in seconds
	}
}

func readFile(cfg *Config) {
	f, err := os.Open(*ConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Reading config: %v\n", *ConfigFilePath)
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatal(err)
	}
}