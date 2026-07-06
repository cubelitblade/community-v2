package port

import (
	"context"
	"net/netip"

	"github.com/cubelitblade/community-v2/backend/services/account/internal/domain/model"
)

// Authenticator authenticates a user by username and password.
type Authenticator interface {
	Authenticate(ctx context.Context, username, password string) (*model.AuthenticatedAccount, error)
}

// LoginRecorder records a successful login event.
type LoginRecorder interface {
	Record(ctx context.Context, accountID int64, ip netip.Addr) error
}
