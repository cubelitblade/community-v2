package bootstrap

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/fx"

	"github.com/cubelitblade/community-v2/backend/services/account/internal/shared"
)

// LoadConfig reads and validates the account service configuration from the
// TOML file located at services/account/config.toml.
func LoadConfig() (*shared.Config, error) {
	loader := koanf.New(".")

	if tomlPath := os.Getenv("ACCOUNT_CONFIG_TOML_PATH"); tomlPath != "" {
		if err := loader.Load(file.Provider(tomlPath), toml.Parser()); err != nil {
			return nil, fmt.Errorf("load config from toml: %w", err)
		}
	}

	if err := loader.Load(env.Provider(".", env.Opt{
		Prefix: "ACCOUNT_",
		TransformFunc: func(key, val string) (string, any) {
			key = strings.ReplaceAll(strings.ToLower(key), "_", ".")
			if strings.Contains(val, " ") {
				return key, strings.Split(val, " ")
			}

			return key, val
		},
	}), nil); err != nil {
		return nil, fmt.Errorf("load config from environment variables: %w", err)
	}

	var cfg shared.Config
	if err := loader.Unmarshal("account", &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

// ConfigModule provides an fx.Option that supplies all sub-configurations as individually injectable dependencies.
func ConfigModule() fx.Option {
	return fx.Options(
		fx.Provide(func(cfg *shared.Config) *shared.DatabaseConfig { return cfg.Database }),
		fx.Provide(func(cfg *shared.Config) *shared.GormConfig { return cfg.Gorm }),
		fx.Provide(func(cfg *shared.Config) *shared.GinConfig { return cfg.Gin }),
		fx.Provide(func(cfg *shared.Config) *shared.SlogConfig { return cfg.Slog }),
		fx.Provide(func(cfg *shared.Config) *shared.SnowflakeConfig { return cfg.Snowflake }),
		fx.Provide(func(cfg *shared.Config) *shared.RabbitMQConfig { return cfg.RabbitMQ }),
		fx.Provide(func(cfg *shared.Config) *shared.JWTConfig { return cfg.JWT }),
		fx.Provide(func(cfg *shared.Config) *shared.CookieConfig { return cfg.Cookie }),
		fx.Provide(func(cfg *shared.Config) *shared.TokenConfig { return cfg.Token }),
	)
}
