package bootstrap

import (
	"fmt"
	"log/slog"
	"time"

	"go.uber.org/fx"

	slogx "github.com/cubelitblade/community-v2/backend/pkg/slog"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/shared"
)

func provideSlogLogger(cfg shared.SlogConfig) (*slog.Logger, error) {
	level, err := slogx.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("parse slog level")
	}

	format, err := slogx.ParseTimeFormat(cfg.Tint.TimeFormat)
	if err != nil {
		format = time.StampMilli
	}

	handler, err := slogx.NewHandler(cfg.Handler, &slog.HandlerOptions{
		Level: level,
	}, slogx.WithTimeFormat(format))
	if err != nil {
		return nil, fmt.Errorf("create handler")
	}

	return slog.New(handler), nil
}

// SlogModule provides the slog fx module.
func SlogModule() fx.Option {
	return fx.Options(
		fx.Provide(provideSlogLogger),
	)
}
