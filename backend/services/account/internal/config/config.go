// Package config holds typed configuration for the account service.
package config

// Config represents the top-level configuration for the account service.
type Config struct {
	ServiceName string `koanf:"service_name"`
	HTTPAddr    string `koanf:"http_addr"`

	Database  *DatabaseConfig  `koanf:"database"`
	Cookie    *CookieConfig    `koanf:"cookie"`
	Redis     *RedisConfig     `koanf:"redis"`
	RabbitMQ  *RabbitMQConfig  `koanf:"rabbitmq"`
	JWT       *JWTConfig       `koanf:"jwt"`
	Snowflake *SnowflakeConfig `koanf:"snowflake"`
	Token     *TokenConfig     `koanf:"token"`
	Slog      *SlogConfig      `koanf:"slog"`
	Gorm      *GormConfig      `koanf:"gorm"`
	Gin       *GinConfig       `koanf:"gin"`
}
