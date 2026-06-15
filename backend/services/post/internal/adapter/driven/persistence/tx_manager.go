package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/cubelitblade/community-v2/backend/services/post/internal/adapter/driven/persistence/ent"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/domain/port"
)

type EntTransactor struct {
	client *ent.Client
}

func NewEntTransactor(client *ent.Client) port.Transactor {
	return &EntTransactor{
		client: client,
	}
}

func (m *EntTransactor) InTx(ctx context.Context, fn func(context.Context) error) error {
	tx, err := m.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("starting tx: %w", err)
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()

	txCtx := ent.NewTxContext(ctx, tx)
	if err := fn(txCtx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return errors.Join(err, rerr)
		}

		return err
	}

	return tx.Commit()
}
