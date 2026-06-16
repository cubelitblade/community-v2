package bootstrap

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"
	"go.uber.org/fx"

	"github.com/cubelitblade/community-v2/backend/services/post/internal/shared"
)

func provideTracerProvider(
	cfg shared.OTELConfig, lc fx.Lifecycle, logger *slog.Logger,
) (*sdktrace.TracerProvider, error) {
	if !cfg.Enable {
		logger.Info("OpenTelemetry is disabled")
		tp := sdktrace.NewTracerProvider()
		otel.SetTracerProvider(tp)

		return tp, nil
	}

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
	}

	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	var exporter sdktrace.SpanExporter
	var err error

	if cfg.ConsoleExporter {
		logger.Info("OpenTelemetry is using stdout trace exporter")

		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout trace exporter: %w", err)
		}
	} else {
		logger.Info("OpenTelemetry is using OTLP trace exporter")

		exporter, err = otlptracegrpc.New(context.Background(), opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
		}
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("shutting down TracerProvider")
			return tp.Shutdown(ctx)
		},
	})

	logger.Info("OpenTelemetry TracerProvider initialized", slog.String("endpoint", cfg.Endpoint))
	return tp, nil
}

func TelemetryModule() fx.Option {
	return fx.Options(
		fx.Provide(provideTracerProvider),
	)
}
