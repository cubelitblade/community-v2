package storage

import (
	"context"
	"fmt"

	"github.com/cubelitblade/community-v2/backend/services/post/internal/post"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/post/storage/ent"
	entPost "github.com/cubelitblade/community-v2/backend/services/post/internal/post/storage/ent/post"
)

type PostWriter struct{}

func (w *PostWriter) Save(ctx context.Context, client *ent.Client, val post.Post) error {
	_, err := client.Post.Create().
		SetID(val.ID()).
		SetAuthorID(val.AuthorID()).
		SetNillableTitle(val.Title()).
		SetContent(val.Content()).
		SetStatus(entPost.Status(val.Status())).
		SetCreatedAt(val.CreatedAt()).
		SetUpdatedAt(val.UpdatedAt()).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("publish post: %w", err)
	}

	return nil
}
