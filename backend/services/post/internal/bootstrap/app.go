package bootstrap

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(LoadConfig),

		DBModule(),
		SlogModule(),
		GRPCModule(),
		SnowflakeModule(),
		TelemetryModule(),

		PostModule(),

		fx.Invoke(func(*grpc.Server) {}),
	)
}
