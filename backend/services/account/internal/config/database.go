package config

// DatabaseConfig holds the configuration for the database connection.
type DatabaseConfig struct {
	URL string `koanf:"url"`
}
