package storage

import (
	"context"
	"fmt"

	"github.com/cubelitblade/community-v2/backend/services/post/internal/post"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/post/storage/ent"
)

type AccountProfileReader struct{}

func (r *AccountProfileReader) Get(ctx context.Context, client *ent.Client, id int64) (*post.AccountProfile, error) {
	row, err := client.AccountProfile.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get account profile: %w", err)
	}

	return &post.AccountProfile{
		ID:          row.ID,
		Username:    row.Username,
		DisplayName: row.DisplayName,
	}, nil
}
