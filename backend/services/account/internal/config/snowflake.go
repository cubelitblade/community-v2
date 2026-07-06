package config

// SnowflakeConfig holds the configuration for the Snowflake distributed ID generator.
type SnowflakeConfig struct {
	NodeID int64 `koanf:"node_id"`
}
