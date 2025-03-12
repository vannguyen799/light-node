package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ProtocolId string `yaml:"protocol-id"`

	LayerEdgeCosmos struct {
		ChainID                   string `yaml:"chainId"`
		RpcEndpoint               string `yaml:"rpcEndpoint"`
		GrpcEndpoint              string `yaml:"grpcEndpoint"`
		AccountPrefix             string `yaml:"accountPrefix"`
		MerkleTreeContractAddress string `yaml:"merkle-tree-contract-address"`
	} `yaml:"layer-edge-cosmos"`

	ZKComputationNode struct {
		NODE_URL string `yaml:"node-url"`
	} `yaml:"zk-computation-node"`
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

	return cfg
}

func validateConfig(cfg *Config) {
	if cfg.ProtocolId == "" {
		log.Fatal("Protocol Id is required in config file")
	}

	if cfg.LayerEdgeCosmos.GrpcEndpoint == "" {
		log.Fatal("LayerEdge Cosmos chain GrpcEndpoint is required")
	}

	if cfg.LayerEdgeCosmos.RpcEndpoint == "" {
		log.Fatal("LayerEdge Cosmos chain RpcEndpoint is required")
	}

	if cfg.LayerEdgeCosmos.MerkleTreeContractAddress == "" {
		log.Fatal("LayerEdge Cosmos chain Merkle tree contract address is required")
	}

	if cfg.ZKComputationNode.NODE_URL == "" {
		log.Fatal("ZK Computation node url is required")
	}
}

func readFile(cfg *Config) {
	var f *os.File
	var err error

	f, err = os.Open(*ConfigFilePath)
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

	validateConfig(cfg)
}
