package bootstrap

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	toml "github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/fx"

	"github.com/cubelitblade/community-v2/backend/services/account/internal/config"
)

// LoadConfig reads and validates the account service configuration from the
// TOML file located at services/account/config.toml.
func LoadConfig() (*config.Config, error) {
	loader := koanf.New(".")

	if err := loader.Load(file.Provider("services/account/config.toml"), toml.Parser()); err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	var cfg config.Config
	if err := loader.Unmarshal("account", &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

// ConfigModule provides an fx.Option that supplies all sub-configurations as individually injectable dependencies.
func ConfigModule() fx.Option {
	return fx.Options(
		fx.Provide(func(cfg *config.Config) *config.DatabaseConfig { return cfg.Database }),
		fx.Provide(func(cfg *config.Config) *config.GormConfig { return cfg.Gorm }),
		fx.Provide(func(cfg *config.Config) *config.GinConfig { return cfg.Gin }),
		fx.Provide(func(cfg *config.Config) *config.SlogConfig { return cfg.Slog }),
		fx.Provide(func(cfg *config.Config) *config.SnowflakeConfig { return cfg.Snowflake }),
		fx.Provide(func(cfg *config.Config) *config.RabbitMQConfig { return cfg.RabbitMQ }),
		fx.Provide(func(cfg *config.Config) *config.JWTConfig { return cfg.JWT }),
		fx.Provide(func(cfg *config.Config) *config.CookieConfig { return cfg.Cookie }),
		fx.Provide(func(cfg *config.Config) *config.TokenConfig { return cfg.Token }),
	)
}
