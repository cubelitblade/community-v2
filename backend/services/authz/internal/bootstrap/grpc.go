package bootstrap

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"

	"github.com/cubelitblade/community-v2/backend/pkg/platform"
)

type grpcParams struct {
	fx.In

	Registrars []platform.GRPCRegistrar `group:"grpc_registrar"`
	Logger     *slog.Logger
}

func provideGRPCServer(lc fx.Lifecycle, p grpcParams, logger *slog.Logger) (*grpc.Server, error) {
	srv := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	for _, r := range p.Registrars {
		r.RegisterGRPC(srv)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("failed to listen for gRPC", "error", err)
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
					log.Printf("GRPC server error: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.GracefulStop()
			return nil
		},
	})

	return srv, nil
}

// func NewGRPCServer(p grpcParams) *grpc.Server {
// 	s := grpc.NewServer()

// 	for _, r := range p.Registrars {
// 		r.RegisterGRPC(s)
// 	}

// 	return s
// }

// func StartGRPC(s *grpc.Server, logger *slog.Logger) {
// 	lis, err := net.Listen("tcp", ":50051")
// 	if err != nil {
// 		logger.Error("failed to listen for gRPC", "error", err)
// 		return
// 	}

// 	logger.Info("Starting gRPC server on :50051")
// 	go s.Serve(lis)
// }

func GRPCModule() fx.Option {
	return fx.Options(
		fx.Provide(provideGRPCServer),
	)
}
