package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	JWTSecret      string
	StaticDir      string
	DBDriver       string
	DBPath         string
	MySQLDSN       string
	RedisEnabled   bool
	RedisAddr      string
	RedisPassword  string
	RedisDB        int
	RedisKeyPrefix string
	CacheTTL       int
	FabricEnabled  bool
	FabricMSP      string
	FabricCertPath string
	FabricKeyPath  string
	FabricTLSPath  string
	FabricPeer     string
	FabricPeerHost string
	FabricChannel  string
	FabricCC       string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret-change-me"),
		StaticDir:      getEnv("STATIC_DIR", ""),
		DBDriver:       getEnv("DB_DRIVER", "sqlite"),
		DBPath:         getEnv("DB_PATH", "trace.db"),
		MySQLDSN:       getEnv("MYSQL_DSN", "root:root@tcp(127.0.0.1:3306)/sea_cucumber_trace?charset=utf8mb4&parseTime=True&loc=Local"),
		RedisEnabled:   getEnvBool("REDIS_ENABLED", false),
		RedisAddr:      getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		RedisDB:        getEnvInt("REDIS_DB", 0),
		RedisKeyPrefix: getEnv("REDIS_KEY_PREFIX", "seatrace:"),
		CacheTTL:       getEnvInt("CACHE_TTL_SECONDS", 300),
		FabricEnabled:  getEnvBool("FABRIC_ENABLED", false),
		FabricMSP:      getEnv("FABRIC_MSP_ID", "Org1MSP"),
		FabricCertPath: getEnv("FABRIC_CERT_PATH", ""),
		FabricKeyPath:  getEnv("FABRIC_KEY_PATH", ""),
		FabricTLSPath:  getEnv("FABRIC_TLS_PATH", ""),
		FabricPeer:     getEnv("FABRIC_PEER_ENDPOINT", "localhost:7051"),
		FabricPeerHost: getEnv("FABRIC_PEER_HOST_ALIAS", ""),
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

func getEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
