package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort           string
	BlockScountClient BlockScountClient
}

type BlockScountClient struct {
	Name                  string
	APIKey                string
	BaseURL               string
	ExternalTxnEndpoint   string
	InternalTxnEndpoint   string
	TokenTransferEndpoint string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	c.AppPort = os.Getenv("APP_PORT")
	c.BlockScountClient.APIKey = os.Getenv("BLOCKSCOUT_API_KEY")
	c.BlockScountClient.BaseURL = os.Getenv("BLOCKSCOUT_BASE_URL")
	c.BlockScountClient.ExternalTxnEndpoint = os.Getenv("BLOCKSCOUT_EXTERNAL_TXN_ENDPOINT")
	c.BlockScountClient.InternalTxnEndpoint = os.Getenv("BLOCKSCOUT_INTERNAL_TXN_ENDPOINT")
	c.BlockScountClient.TokenTransferEndpoint = os.Getenv("BLOCKSCOUT_TOKEN_TRANSFER_ENDPOINT")
}
