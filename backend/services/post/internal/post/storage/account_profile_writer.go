package storage

import (
	"context"
	"fmt"

	"github.com/cubelitblade/community-v2/backend/services/post/internal/post"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/post/storage/ent"
)

type AccountProfileWriter struct{}

func (w *AccountProfileWriter) Save(ctx context.Context, client *ent.Client, val post.AccountProfile) error {
	_, err := client.AccountProfile.Create().
		SetID(val.ID).
		SetUsername(val.Username).
		SetDisplayName(val.DisplayName).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("write account profile: %w", err)
	}

	return nil
}
