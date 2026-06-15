package application

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cubelitblade/community-v2/backend/services/post/internal/domain"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/domain/model"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/domain/port"
)

type PublishPost struct {
	postSaver  port.PostSaver
	authz      port.Authorizer
	transactor port.Transactor
	idgen      port.IDGenerator

	logger   *slog.Logger
	timeFunc func() time.Time
}

func NewPublishPost(
	postSaver port.PostSaver, authz port.Authorizer, transactor port.Transactor, idgen port.IDGenerator,
	logger *slog.Logger,
) *PublishPost {
	return &PublishPost{
		postSaver:  postSaver,
		authz:      authz,
		transactor: transactor,
		logger:     logger,
		timeFunc:   time.Now,
	}
}

func (p *PublishPost) Publish(ctx context.Context, accountID int64, title *string, content string) error {
	ok, err := p.authz.CanPublishPost(ctx, accountID)
	if err != nil {
		p.logger.Error("failed to check permission", slog.Any("error", err))
	}
	if !ok {
		p.logger.Warn("permission denied",
			slog.Int64("accountID", accountID),
			slog.String("permission", "can_publish_post"),
		)
		return domain.ErrPermissionDenied
	}
	p.logger.Debug("permission check passed", slog.Int64("account_id", accountID))

	id, err := p.idgen.Next()
	if err != nil {
		p.logger.Error("failed to generate id", slog.Any("error", err))
		return fmt.Errorf("generate id: %w", err)
	}

	post, err := model.NewPost(id, accountID, title, content, p.timeFunc())
	if err != nil {
		p.logger.Error("failed to create post", slog.Any("error", err))
		return fmt.Errorf("create post: %w", err)
	}

	if err := p.transactor.InTx(ctx, func(ctx context.Context) error {
		if err := p.postSaver.Save(ctx, *post); err != nil {
			return err
		}

		// TODO: Write outbox.

		return nil
	}); err != nil {
		p.logger.Error("transaction error", slog.Any("error", err))
		return fmt.Errorf("tx error: %w", err)
	}
	p.logger.Info("post published successfully",
		slog.Int64("postID", post.ID()),
		slog.Int64("accountID", accountID),
	)

	return nil
}
