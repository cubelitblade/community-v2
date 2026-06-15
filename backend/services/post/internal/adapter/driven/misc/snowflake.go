package misc

import (
	"github.com/cubelitblade/community-v2/backend/pkg/common/idgen"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/domain/port"
)

type IDGenerator struct {
	idgen *idgen.Snowflake
}

func NewIDGenerator(idgen *idgen.Snowflake) port.IDGenerator {
	return &IDGenerator{
		idgen: idgen,
	}
}

func (g *IDGenerator) Next() (int64, error) {
	return g.idgen.NextID()
}
