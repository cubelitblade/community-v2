//nolint:testpackage // White-box domain tests
package authz

import (
	"context"
	"errors"
	"testing"

	openfga "github.com/openfga/go-sdk"
	gosdk "github.com/openfga/go-sdk/client"
)

var errFGA = errors.New("openfga error")

type stubCheckRequest struct {
	allowed *bool
	err     error
}

//nolint:ireturn // returning interface required by SDK interface
func (s *stubCheckRequest) Options(_ gosdk.ClientCheckOptions) gosdk.SdkClientCheckRequestInterface {
	return s
}

//nolint:ireturn // returning interface required by SDK interface
func (s *stubCheckRequest) Body(_ gosdk.ClientCheckRequest) gosdk.SdkClientCheckRequestInterface {
	return s
}

func (s *stubCheckRequest) Execute() (*gosdk.ClientCheckResponse, error) {
	if s.err != nil {
		return nil, s.err
	}

	return &gosdk.ClientCheckResponse{CheckResponse: openfga.CheckResponse{Allowed: s.allowed}}, nil
}

//nolint:revive // Interface method name from OpenFGA SDK
func (s *stubCheckRequest) GetAuthorizationModelIdOverride() *string { return nil }

//nolint:revive // Interface method name from OpenFGA SDK
func (s *stubCheckRequest) GetStoreIdOverride() *string { return nil }

func (s *stubCheckRequest) GetContext() context.Context { return context.Background() }

func (s *stubCheckRequest) GetBody() *gosdk.ClientCheckRequest { return nil }

func (s *stubCheckRequest) GetOptions() *gosdk.ClientCheckOptions { return nil }

func newStubCheckRequest(allowed *bool, err error) *stubCheckRequest {
	return &stubCheckRequest{allowed: allowed, err: err}
}

type stubGetStoreRequest struct {
	err error
}

//nolint:ireturn // returning interface required by SDK interface
func (s *stubGetStoreRequest) Options(_ gosdk.ClientGetStoreOptions) gosdk.SdkClientGetStoreRequestInterface {
	return s
}

func (s *stubGetStoreRequest) Execute() (*gosdk.ClientGetStoreResponse, error) {
	if s.err != nil {
		return nil, s.err
	}

	return &gosdk.ClientGetStoreResponse{}, nil
}

//nolint:revive // Interface method name from OpenFGA SDK
func (s *stubGetStoreRequest) GetAuthorizationModelIdOverride() *string { return nil }

//nolint:revive // Interface method name from OpenFGA SDK
func (s *stubGetStoreRequest) GetStoreIdOverride() *string { return nil }

func (s *stubGetStoreRequest) GetContext() context.Context { return context.Background() }

func (s *stubGetStoreRequest) GetOptions() *gosdk.ClientGetStoreOptions { return nil }

func newStubGetStoreRequest(err error) *stubGetStoreRequest {
	return &stubGetStoreRequest{err: err}
}

type stubFGAClient struct {
	checkResult    *stubCheckRequest
	getStoreResult *stubGetStoreRequest
}

//nolint:ireturn // returning interface required by SDK interface
func (c *stubFGAClient) Check(_ context.Context) gosdk.SdkClientCheckRequestInterface {
	return c.checkResult
}

//nolint:ireturn // returning interface required by SDK interface
func (c *stubFGAClient) GetStore(_ context.Context) gosdk.SdkClientGetStoreRequestInterface {
	return c.getStoreResult
}

func newStubFGAClient(check *stubCheckRequest, getStore *stubGetStoreRequest) *stubFGAClient {
	return &stubFGAClient{checkResult: check, getStoreResult: getStore}
}

func TestChecker_CheckAllowed(t *testing.T) {
	t.Parallel()

	checker := NewChecker(newStubFGAClient(
		newStubCheckRequest(openfga.PtrBool(true), nil),
		newStubGetStoreRequest(nil),
	))

	allowed, err := checker.Check(context.Background(), "user:alice", "can_edit", "post:1")
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}

	if !allowed {
		t.Error("Check() = false, want true")
	}
}

func TestChecker_CheckDenied(t *testing.T) {
	t.Parallel()

	checker := NewChecker(newStubFGAClient(
		newStubCheckRequest(openfga.PtrBool(false), nil),
		newStubGetStoreRequest(nil),
	))

	allowed, err := checker.Check(context.Background(), "user:bob", "can_edit", "post:1")
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}

	if allowed {
		t.Error("Check() = true, want false")
	}
}

func TestChecker_CheckNilAllowed(t *testing.T) {
	t.Parallel()

	checker := NewChecker(newStubFGAClient(
		newStubCheckRequest(nil, nil),
		newStubGetStoreRequest(nil),
	))

	allowed, err := checker.Check(context.Background(), "user:carol", "can_view", "post:1")
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}

	if allowed {
		t.Error("Check() = true, want false for nil Allowed")
	}
}

func TestChecker_CheckError(t *testing.T) {
	t.Parallel()

	checker := NewChecker(newStubFGAClient(
		newStubCheckRequest(nil, errFGA),
		newStubGetStoreRequest(nil),
	))

	_, err := checker.Check(context.Background(), "user:alice", "can_edit", "post:1")
	if err == nil {
		t.Fatal("Check() expected error, got nil")
	}
}

func TestChecker_PingOK(t *testing.T) {
	t.Parallel()

	checker := NewChecker(newStubFGAClient(
		newStubCheckRequest(nil, nil),
		newStubGetStoreRequest(nil),
	))

	err := checker.Ping(context.Background())
	if err != nil {
		t.Fatalf("Ping() error = %v", err)
	}
}

func TestChecker_PingError(t *testing.T) {
	t.Parallel()

	checker := NewChecker(newStubFGAClient(
		newStubCheckRequest(nil, nil),
		newStubGetStoreRequest(errFGA),
	))

	err := checker.Ping(context.Background())
	if err == nil {
		t.Fatal("Ping() expected error, got nil")
	}
}
