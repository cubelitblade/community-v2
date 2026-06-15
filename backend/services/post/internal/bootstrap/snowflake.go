package bootstrap

import (
	"fmt"

	"go.uber.org/fx"

	"github.com/cubelitblade/community-v2/backend/pkg/common/idgen"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/shared"
)

func provideSnowflake(cfg shared.SnowflakeConfig) (*idgen.Snowflake, error) {
	ids, err := idgen.NewSnowflake(cfg.NodeID)
	if err != nil {
		return nil, fmt.Errorf("generate snowflake instance: %w", err)
	}

	return ids, nil
}

func SnowflakeModule() fx.Option {
	return fx.Options(
		fx.Provide(provideSnowflake),
	)
}
