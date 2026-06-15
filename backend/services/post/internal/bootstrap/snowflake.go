package bootstrap

import (
	"github.com/cubelitblade/community-v2/backend/pkg/common/idgen"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/shared"
	"go.uber.org/fx"
)

func provideSnowflake(cfg shared.SnowflakeConfig) (*idgen.Snowflake, error) {
	return idgen.NewSnowflake(cfg.NodeID)
}

func SnowflakeModule() fx.Option {
	return fx.Options(
		fx.Provide(provideSnowflake),
	)
}
