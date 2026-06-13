package ent_test

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"github.com/cubelitblade/community-v2/backend/services/post/internal/post/storage/ent/enttest"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/post/storage/ent/post"
)

var (
	validContent = "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
	validTitle   = "Lorem Ipsum"
)

func TestPost_Create(t *testing.T) {
	t.Parallel()

	dsn := fmt.Sprintf("file:%s?mode=memory&_fk=1", t.Name())

	client := enttest.Open(t, "sqlite3", dsn)
	defer client.Close()

	ctx := context.Background()

	author, err := client.AccountProfile.Create().
		SetUsername("test_author").
		SetDisplayName("Test Author").
		Save(ctx)
	if err != nil {
		t.Fatalf("could not create prerequisite author: %v", err)
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		authorID   int64
		title      *string
		content    string
		status     post.Status
		wantErr    bool
		wantTitle  *string
		wantStatus post.Status
	}{
		{
			name:       "happy_path/default_values",
			authorID:   author.ID,
			title:      nil,
			content:    validContent,
			status:     "",
			wantErr:    false,
			wantTitle:  nil,
			wantStatus: "published",
		},
		{
			name:       "happy_path/with_title",
			authorID:   author.ID,
			title:      &validTitle,
			content:    validContent,
			status:     "",
			wantErr:    false,
			wantTitle:  &validTitle,
			wantStatus: "published",
		},
		{
			name:      "validate/invalid_status_enum",
			authorID:  author.ID,
			title:     nil,
			content:   validContent,
			status:    "deleted",
			wantErr:   true,
			wantTitle: nil,
		},
		{
			name:      "validate/missing_required_content",
			authorID:  author.ID,
			title:     nil,
			content:   "",
			status:    "",
			wantErr:   true,
			wantTitle: nil,
		},
		{
			name:      "rule/non_existent_author_id",
			authorID:  99999,
			title:     nil,
			content:   validContent,
			status:    "",
			wantErr:   true,
			wantTitle: nil,
		},
	}
	//nolint:paralleltest // Serial execution to avoid SQLite concurrent write lock conflicts.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := client.Post.Create().
				SetContent(tt.content)

			if tt.authorID != 0 {
				builder.SetAuthorID(tt.authorID)
			}

			if tt.title != nil {
				builder.SetTitle(*tt.title)
			}

			if tt.status != "" {
				builder.SetStatus(tt.status)
			}

			got, gotErr := builder.Save(ctx)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			if gotErr != nil {
				t.Fatalf("CreatePost() succeeded unexpectedly: %v", gotErr)
			}

			assert.NotZero(t, got.ID, "ID should be auto-generated")
			assert.Equal(t, tt.wantTitle, got.Title)
			assert.Equal(t, tt.wantStatus, got.Status)
			assert.Equal(t, tt.content, got.Content)
			assert.NotZero(t, got.CreatedAt, "CreatedAt should have default value")
			assert.NotZero(t, got.UpdatedAt, "UpdatedAt should have default value")
		})
	}
}

func TestAccountProfile_Delete(t *testing.T) {
	t.Parallel()

	dsn := fmt.Sprintf("file:%s?mode=memory&_fk=1", t.Name())

	client := enttest.Open(t, "sqlite3", dsn)
	defer client.Close()

	ctx := context.Background()

	tests := []struct {
		name          string
		hasLinkedPost bool
		wantErr       bool
	}{
		{
			name:          "happy_path/delete_author_without_posts",
			hasLinkedPost: false,
			wantErr:       false,
		},
		{
			name:          "rule/on_delete_no_action_blocks_deletion",
			hasLinkedPost: true,
			wantErr:       true,
		},
	}
	//nolint:paralleltest // Serial execution to avoid SQLite concurrent write lock conflicts.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			author, err := client.AccountProfile.Create().
				SetUsername("temp_author").
				SetDisplayName("Temp").
				Save(ctx)
			if err != nil {
				t.Fatalf("could not create prerequisite author: %v", err)
			}

			if tt.hasLinkedPost {
				_, err = client.Post.Create().
					SetContent("linked content").
					SetAuthorID(author.ID).
					Save(ctx)
				if err != nil {
					t.Fatalf("could not create linked post: %v", err)
				}
			}

			gotErr := client.AccountProfile.DeleteOneID(author.ID).Exec(ctx)

			if tt.wantErr {
				assert.Error(t, gotErr, "deleting author with linked posts should be blocked by NoAction")
				return
			}

			if gotErr != nil {
				t.Fatalf("DeleteAuthor() failed unexpectedly: %v", gotErr)
			}
		})
	}
}
