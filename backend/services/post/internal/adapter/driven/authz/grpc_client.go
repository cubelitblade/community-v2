package authz

import (
	"context"
	"fmt"

	authzv1 "github.com/cubelitblade/community-v2/backend/gen/go/authz/v1"
	"github.com/cubelitblade/community-v2/backend/services/authz/contract"
	"github.com/cubelitblade/community-v2/backend/services/post/internal/domain/port"
)

type Authorizer struct {
	client authzv1.AuthzServiceClient
}

func NewAuthorizer(client authzv1.AuthzServiceClient) port.Authorizer {
	return &Authorizer{client: client}
}

func (a *Authorizer) CanPublishPost(ctx context.Context, accountID int64) (bool, error) {
	resp, err := a.client.Check(ctx, &authzv1.CheckRequest{
		User:     fmt.Sprintf("user:%d", accountID),
		Relation: contract.RelationCanPublishPost,
		Object:   contract.ObjectSystemCommunity,
	})
	if err != nil {
		return false, err
	}
	return resp.Allowed, nil
}
