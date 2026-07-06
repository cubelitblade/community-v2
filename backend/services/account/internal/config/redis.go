package config

// RedisConfig holds the configuration for the Redis client connection.
type RedisConfig struct {
	Addr string `koanf:"addr"`
}
