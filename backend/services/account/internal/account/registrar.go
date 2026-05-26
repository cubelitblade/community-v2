package account

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/cubelitblade/community-v2/backend/pkg/common/idgen"
)

type creator interface {
	Create(ctx context.Context, acc *Account) error
}

type eventPublisher interface {
	Publish(ctx context.Context, routingKey string, body []byte) error
}

// Registrar orchestrates the account registration process.
type Registrar struct {
	ids       idgen.Generator
	creator   creator
	publisher eventPublisher
	now       func() time.Time
	logger    *slog.Logger
}

// NewRegistrar creates and returns a new Registrar.
func NewRegistrar(
	ids idgen.Generator,
	creator creator,
	publisher eventPublisher,
	logger *slog.Logger,
) *Registrar {
	if logger == nil {
		panic("nil logger")
	}

	return &Registrar{
		ids:       ids,
		creator:   creator,
		publisher: publisher,
		now:       time.Now,
		logger:    logger,
	}
}

// Register creates a new account, persists it, and returns the result.
func (r *Registrar) Register(
	ctx context.Context,
	username, password string,
) (Account, error) {
	id, err := r.ids.NextID()
	if err != nil {
		return Account{}, fmt.Errorf(
			"could not generate account ID: %w",
			err,
		)
	}

	acc, err := Register(id, username, password, r.now())
	if err != nil {
		return Account{}, fmt.Errorf(
			"could not register account: %w",
			err,
		)
	}

	if err := r.creator.Create(ctx, &acc); err != nil {
		return Account{}, fmt.Errorf(
			"could not create account: %w",
			err,
		)
	}

	r.logger.Info(
		"account created",
		"account_id", acc.id,
		"username", acc.username.Value(),
		"display_name", acc.displayName,
		"role", acc.role,
		"status", acc.status,
		"created_at", acc.createdAt.Format(time.RFC3339),
	)

	if r.publisher != nil {
		payload, err := json.Marshal(userRegisteredPayload{UserID: strconv.FormatInt(int64(acc.id), 10)})
		if err != nil {
			r.logger.Error("failed to marshal user.registered event", "error", err)
		} else if err := r.publisher.Publish(ctx, "user.registered", payload); err != nil {
			r.logger.Error("failed to publish user.registered event", "error", err)
		}
	}

	return acc, nil
}

// userRegisteredPayload is the JSON payload for the user.registered event.

//nolint:tagliatelle // JSON keys are cross-service contract
type userRegisteredPayload struct {
	UserID string `json:"user_id"`
}
