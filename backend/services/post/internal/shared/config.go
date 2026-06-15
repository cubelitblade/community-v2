// Package shared contains common types and configurations used across different
// layers of the post service. Currently, it primarily defines the configuration
// structures required to bootstrap the application and initialize adapters.
package shared

import (
	"go.uber.org/fx"
)

type AppEnv string

type Config struct {
	fx.Out

	Env       AppEnv          `koanf:"env" validate:"required,oneof=dev prod"`
	Gin       GinConfig       `koanf:"gin"`
	Database  DatabaseConfig  `koanf:"db"`
	Slog      SlogConfig      `koanf:"slog"`
	GRPC      GRPCConfig      `koanf:"grpc"`
	Snowflake SnowflakeConfig `koanf:"snowflake"`
}

type GinConfig struct {
	Mode string `koanf:"mode" validate:"required,oneof=debug test release"`
}

type DatabaseConfig struct {
	URI string `koanf:"uri" validate:"required,uri"`
}

type SlogConfig struct {
	Handler string         `koanf:"handler" validate:"oneof=text json tint default discard"`
	Level   string         `koanf:"level" validate:"oneof=debug info warn error"`
	Tint    SlogTintConfig `koanf:"tint"`
}

type SlogTintConfig struct {
	TimeFormat string `koanf:"time_format"`
}

type GRPCConfig struct {
	Reflection bool `koanf:"reflection"`
}

type SnowflakeConfig struct {
	NodeID int64 `koanf:"node_id"`
}
