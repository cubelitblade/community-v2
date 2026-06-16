package grpc

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	postv1 "github.com/cubelitblade/community-v2/backend/gen/go/post/v1"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/application"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/domain"
)

var _ postv1.PostServiceServer = (*PostServiceServer)(nil)

var errorMap = map[error]error{
	domain.ErrContentBlank:     status.Error(codes.InvalidArgument, "content is blank"),
	domain.ErrIDInvalid:        status.Error(codes.InvalidArgument, "invalid id"),
	domain.ErrPermissionDenied: status.Error(codes.PermissionDenied, "permission denied"),
}

type PostServiceServer struct {
	postv1.UnimplementedPostServiceServer

	publishPost *application.PublishPost
}

func NewPostServiceServer(publishPost *application.PublishPost) *PostServiceServer {
	return &PostServiceServer{
		UnimplementedPostServiceServer: postv1.UnimplementedPostServiceServer{},
		publishPost:                    publishPost,
	}
}

func (s *PostServiceServer) Publish(ctx context.Context, req *postv1.PublishRequest) (*postv1.PublishResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.Int64("account.id", req.AccountId),
	)

	id, err := s.publishPost.Publish(ctx, req.AccountId, req.Title, req.Content)
	if err != nil {
		gerr, ok := errorMap[err]
		if !ok {
			//nolint:errcheck // Unwrap for gRPC to handle.
			return nil, status.Error(codes.Internal, "unexpected error")
		}

		return nil, gerr
	}

	span.SetAttributes(
		attribute.Int64("post.id", id),
	)

	return &postv1.PublishResponse{
		PostId: id,
	}, nil
}
