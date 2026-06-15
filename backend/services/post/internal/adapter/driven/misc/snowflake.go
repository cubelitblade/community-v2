// Package misc contains miscellaneous, stateless driven adapters.
// It implements various utility ports that do not maintain any application state,
// providing infrastructure support for functional and side-effect-free operations.
package misc

import (
	"fmt"

	"github.com/cubelitblade/community-v2/backend/pkg/common/idgen"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/domain/port"
)

type IDGenerator struct {
	ids *idgen.Snowflake
}

func NewIDGenerator(ids *idgen.Snowflake) port.IDGenerator {
	return &IDGenerator{
		ids: ids,
	}
}

func (g *IDGenerator) Next() (int64, error) {
	id, err := g.ids.NextID()
	if err != nil {
		return 0, fmt.Errorf("generate id: %w", err)
	}

	return id, nil
}
