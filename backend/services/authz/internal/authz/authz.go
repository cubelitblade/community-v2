package authz

import (
	"context"
	"fmt"

	gosdk "github.com/openfga/go-sdk/client"
)

// FGAConfig holds the OpenFGA connection settings.
type FGAConfig struct {
	APIURL  string
	StoreID string
	ModelID string
}

// goSdkClient is the subset of the OpenFGA Go SDK client that Checker needs.
type goSdkClient interface {
	Check(ctx context.Context) gosdk.SdkClientCheckRequestInterface
	GetStore(ctx context.Context) gosdk.SdkClientGetStoreRequestInterface
}

// Checker performs authorization checks against OpenFGA.
type Checker struct {
	client goSdkClient
}

// NewChecker creates a Checker backed by the given OpenFGA SDK client.
func NewChecker(client goSdkClient) *Checker {
	return &Checker{
		client: client,
	}
}

// Check returns whether the given user has the given relation to the given object.
func (c *Checker) Check(ctx context.Context, user, relation, object string) (bool, error) {
	body := gosdk.ClientCheckRequest{
		User:     user,
		Relation: relation,
		Object:   object,
	}

	resp, err := c.client.Check(ctx).Body(body).Execute()
	if err != nil {
		return false, fmt.Errorf("openfga check: %w", err)
	}

	if resp.Allowed == nil {
		return false, nil
	}

	return *resp.Allowed, nil
}

// Ping verifies connectivity to OpenFGA by reading the store.
func (c *Checker) Ping(ctx context.Context) error {
	_, err := c.client.GetStore(ctx).Execute()
	if err != nil {
		return fmt.Errorf("openfga ping: %w", err)
	}

	return nil
}

// Ensure *OpenFgaClient satisfies goSdkClient.
var _ goSdkClient = (*gosdk.OpenFgaClient)(nil)
