package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	JWTSecret      string
	DBPath         string
	FabricEnabled  bool
	FabricMSP      string
	FabricCertPath string
	FabricKeyPath  string
	FabricTLSPath  string
	FabricPeer     string
	FabricChannel  string
	FabricCC       string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret-change-me"),
		DBPath:         getEnv("DB_PATH", "trace.db"),
		FabricEnabled:  getEnvBool("FABRIC_ENABLED", false),
		FabricMSP:      getEnv("FABRIC_MSP_ID", "Org1MSP"),
		FabricCertPath: getEnv("FABRIC_CERT_PATH", ""),
		FabricKeyPath:  getEnv("FABRIC_KEY_PATH", ""),
		FabricTLSPath:  getEnv("FABRIC_TLS_PATH", ""),
		FabricPeer:     getEnv("FABRIC_PEER_ENDPOINT", "localhost:7051"),
		FabricChannel:  getEnv("FABRIC_CHANNEL", "mychannel"),
		FabricCC:       getEnv("FABRIC_CHAINCODE", "traceability"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}
