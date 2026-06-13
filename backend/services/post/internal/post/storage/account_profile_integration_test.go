//go:build integration

package storage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cubelitblade/community-v2/backend/services/post/internal/post"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/post/storage"
)

func TestIntegrationAccountProfileWorkflow(t *testing.T) {
	ctx := context.Background()

	client := newClient(ctx, t)

	writer := storage.AccountProfileWriter{}
	reader := storage.AccountProfileReader{}

	t.Run("", func(t *testing.T) {
		profile := post.AccountProfile{
			ID:          authorID,
			Username:    username,
			DisplayName: displayName,
		}

		err := writer.Save(ctx, client, profile)
		if err != nil {
			t.Fatalf("Save() failed: %v", err)
		}

		got, err := reader.Get(ctx, client, authorID)
		if err != nil {
			t.Fatalf("Get() failed: %v", err)
		}

		assert.Equal(t, profile.ID, got.ID)
		assert.Equal(t, profile.Username, got.Username)
		assert.Equal(t, profile.DisplayName, got.DisplayName)
	})
}
