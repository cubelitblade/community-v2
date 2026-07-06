// Package port declares output ports (driven-adapter interfaces) for the account domain.
package port

import (
	"context"

	"gorm.io/gorm"

	"github.com/cubelitblade/community-v2/backend/pkg/events/outbox"
	"github.com/cubelitblade/community-v2/backend/services/account/internal/domain/model"
)

// Creator persists a new account.
type Creator interface {
	Create(ctx context.Context, db *gorm.DB, acc *model.Account) error
}

// OutboxSaver persists an outbox event for asynchronous processing.
type OutboxSaver interface {
	Save(ctx context.Context, db *gorm.DB, event *outbox.Entry) error
}

// ByIDFinder finds an account by its unique ID.
type ByIDFinder interface {
	FindByID(ctx context.Context, db *gorm.DB, id int64) (*model.Account, error)
}

// ByUsernameFinder finds an account by username.
type ByUsernameFinder interface {
	FindByUsername(ctx context.Context, db *gorm.DB, username string) (*model.Account, error)
}

// LastLoginUpdater updates the last-login timestamp/IP for an account.
type LastLoginUpdater interface {
	UpdateLastLogin(ctx context.Context, db *gorm.DB, acc *model.Account) error
}
