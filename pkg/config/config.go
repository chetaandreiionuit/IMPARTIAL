package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort             string  `mapstructure:"SERVER_PORT"`
	TemporalHost           string  `mapstructure:"TEMPORAL_HOST"`
	DgraphHost             string  `mapstructure:"DGRAPH_HOST"`
	GeminiAPIKey           string  `mapstructure:"GEMINI_API_KEY"`
	NewsAPIKey             string  `mapstructure:"NEWS_API_KEY"` // New field for Global News
	ArweaveKeyPath         string  `mapstructure:"ARWEAVE_KEY_PATH"`
	ArweaveGateway         string  `mapstructure:"ARWEAVE_GATEWAY"`
	DBURL                  string  `mapstructure:"DB_URL"`
	DgraphCertPath         string  `mapstructure:"DGRAPH_CERT_PATH"`
	SolanaPrivateKey       string  `mapstructure:"SOLANA_PRIVATE_KEY"`
	SolanaRPC              string  `mapstructure:"SOLANA_RPC"`
	DeduplicationThreshold float64 `mapstructure:"DEDUPLICATION_THRESHOLD"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("TEMPORAL_HOST", "localhost:7233")
	viper.SetDefault("DGRAPH_HOST", "localhost:9080")
	viper.SetDefault("ARWEAVE_GATEWAY", "https://arweave.net")
	viper.SetDefault("DB_URL", "postgres://user:password@localhost:5432/truthweave?sslmode=disable")
	viper.SetDefault("DEDUPLICATION_THRESHOLD", 0.90)

	viper.AutomaticEnv()

	// [RO] Încearcă să încarci .env
	viper.SetConfigFile(".env")
	viper.ReadInConfig() // Ignorăm eroarea dacă nu există (folosim env vars sau defaults)

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
