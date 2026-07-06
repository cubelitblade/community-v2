package config

// GinConfig holds the configuration for the Gin web framework.
type GinConfig struct {
	Mode string `koanf:"mode"`
}
