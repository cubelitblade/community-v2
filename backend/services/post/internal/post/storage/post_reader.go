package storage

import (
	"context"
	"fmt"

	"github.com/cubelitblade/community-v2/backend/services/post/internal/post"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/post/storage/ent"
)

type PostReader struct{}

func (r *PostReader) Get(ctx context.Context, client *ent.Client, id int64) (*post.Post, error) {
	row, err := client.Post.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get post: %w", err)
	}

	return post.Unmarshall(
		row.ID, row.AuthorID, row.Title, row.Content, post.Status(row.Status), row.CreatedAt, row.UpdatedAt,
	), nil
}
