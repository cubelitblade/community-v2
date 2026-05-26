package bootstrap

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/caarlos0/env"
	pkgEvents "github.com/cubelitblade/community-v2/backend/pkg/events"
	platformconfig "github.com/cubelitblade/community-v2/backend/pkg/platform/config"
	"github.com/cubelitblade/community-v2/backend/services/authz/internal/authz"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Config holds all application configuration grouped by concern.
type Config struct {
	Slog     *platformconfig.SlogConfig
	Gin      *platformconfig.GinConfig
	Server   *platformconfig.ServerConfig
	FGA      *authz.FGAConfig
	RabbitMQ *pkgEvents.Config
}

// ConfigOut provides all config sub-structs as separate fx dependencies.
type ConfigOut struct {
	fx.Out

	Slog     *platformconfig.SlogConfig
	Gin      *platformconfig.GinConfig
	Server   *platformconfig.ServerConfig
	FGA      *authz.FGAConfig
	RabbitMQ *pkgEvents.Config
}

// RawConfig is the flat env-var representation parsed by caarlos0/env.
type RawConfig struct {
	AppEnv      string `env:"APP_ENV" envDefault:"dev"`
	HTTPAddr    string `env:"HTTP_ADDR" envDefault:":8081"`
	FGAAPIURL   string `env:"FGA_API_URL" required:"true"`
	FGAStoreID  string `env:"FGA_STORE_ID" required:"true"`
	FGAModelID  string `env:"FGA_MODEL_ID"`
	RabbitMQURL string `env:"RABBITMQ_URL" required:"true"`
}

// LoadConfig reads environment variables and returns a validated Config.
func LoadConfig() (*Config, error) {
	var raw RawConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &Config{
		Slog:     raw.newSlogConfig(),
		Gin:      raw.newGinConfig(),
		Server:   raw.newServerConfig(),
		FGA:      raw.newFGAConfig(),
		RabbitMQ: raw.newRabbitMQConfig(),
	}, nil
}

// ProvideConfig loads configuration from the environment and exposes each
// config struct as a separate fx dependency via ConfigOut.
func ProvideConfig() (ConfigOut, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return ConfigOut{}, fmt.Errorf("failed to load config: %w", err)
	}

	return ConfigOut{
		Out:      fx.Out{},
		Slog:     cfg.Slog,
		Gin:      cfg.Gin,
		Server:   cfg.Server,
		FGA:      cfg.FGA,
		RabbitMQ: cfg.RabbitMQ,
	}, nil
}

func (r *RawConfig) newSlogConfig() *platformconfig.SlogConfig {
	switch r.AppEnv {
	case "prod":
		return &platformconfig.SlogConfig{
			LogLevel:  slog.LevelInfo,
			LogFormat: platformconfig.LogFormatJSON,
		}
	default:
		return &platformconfig.SlogConfig{
			LogLevel:  slog.LevelDebug,
			LogFormat: platformconfig.LogFormatText,
		}
	}
}

func (r *RawConfig) newGinConfig() *platformconfig.GinConfig {
	switch r.AppEnv {
	case "prod":
		return &platformconfig.GinConfig{
			GinMode:        gin.ReleaseMode,
			TrustedProxies: nil,
		}
	default:
		return &platformconfig.GinConfig{
			GinMode:        gin.DebugMode,
			TrustedProxies: nil,
		}
	}
}

func (r *RawConfig) newServerConfig() *platformconfig.ServerConfig {
	return &platformconfig.ServerConfig{
		HTTPAddr:          r.HTTPAddr,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}
}

func (r *RawConfig) newFGAConfig() *authz.FGAConfig {
	return &authz.FGAConfig{
		APIURL:  r.FGAAPIURL,
		StoreID: r.FGAStoreID,
		ModelID: r.FGAModelID,
	}
}

func (r *RawConfig) newRabbitMQConfig() *pkgEvents.Config {
	return &pkgEvents.Config{
		URL:          r.RabbitMQURL,
		ExchangeName: "domain.events",
		QueueName:    "authz.tuple-writer",
		ConsumerTag:  "authz.tuple-writer",
		Bindings: []string{
			"user.registered",
			"user.role_assigned",
			"user.role_revoked",
			"account.created",
			"post.created",
		},
	}
}
