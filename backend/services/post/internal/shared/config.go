package shared

import (
	"go.uber.org/fx"
)

type AppEnv string

type Config struct {
	fx.Out

	Env      AppEnv         `koanf:"env" validate:"required,oneof=dev prod"`
	Gin      GinConfig      `koanf:"gin"`
	Database DatabaseConfig `koanf:"db"`
	Slog     SlogConfig     `koanf:"slog"`
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
