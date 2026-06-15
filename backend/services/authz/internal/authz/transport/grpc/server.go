package grpc

import (
	"context"

	"google.golang.org/grpc"

	authzv1 "github.com/cubelitblade/community-v2/backend/gen/go/authz/v1"
)

// Checker verifies if a user has a specific relation to an object.
type Checker interface {
	Check(ctx context.Context, user, relation, object string) (bool, error)
}

type Server struct {
	authzv1.UnimplementedAuthzServiceServer

	checker Checker
}

func NewServer(checker Checker) *Server {
	return &Server{
		checker: checker,
	}
}

func (s *Server) Check(ctx context.Context, req *authzv1.CheckRequest) (*authzv1.CheckResponse, error) {
	allowed, err := s.checker.Check(ctx, req.User, req.Relation, req.Object)
	if err != nil {
		return nil, err
	}

	return &authzv1.CheckResponse{
		Allowed: allowed,
	}, nil
}

func (s *Server) RegisterGRPC(srv *grpc.Server) {
	authzv1.RegisterAuthzServiceServer(srv, s)
}
