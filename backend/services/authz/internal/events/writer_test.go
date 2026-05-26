//nolint:testpackage // White-box tests for Writer with unexported fgaWriteClient interface
package events

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	gosdk "github.com/openfga/go-sdk/client"
)

var errFGA = errors.New("openfga error")

var testLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

type stubWriteRequest struct {
	err error
}

//nolint:ireturn // returning interface required by SDK interface
func (s *stubWriteRequest) Options(_ gosdk.ClientWriteOptions) gosdk.SdkClientWriteRequestInterface {
	return s
}

//nolint:ireturn // returning interface required by SDK interface
func (s *stubWriteRequest) Body(_ gosdk.ClientWriteRequest) gosdk.SdkClientWriteRequestInterface {
	return s
}

func (s *stubWriteRequest) Execute() (*gosdk.ClientWriteResponse, error) {
	if s.err != nil {
		return nil, s.err
	}

	return &gosdk.ClientWriteResponse{}, nil
}

//nolint:revive // Interface method name from OpenFGA SDK
func (s *stubWriteRequest) GetAuthorizationModelIdOverride() *string { return nil }

//nolint:revive // Interface method name from OpenFGA SDK
func (s *stubWriteRequest) GetStoreIdOverride() *string { return nil }

func (s *stubWriteRequest) GetContext() context.Context { return context.Background() }

func (s *stubWriteRequest) GetOptions() *gosdk.ClientWriteOptions { return nil }

func (s *stubWriteRequest) GetBody() *gosdk.ClientWriteRequest { return nil }

func newStubWrite(err error) *stubWriteRequest {
	return &stubWriteRequest{err: err}
}

type stubFGAWriteClient struct {
	writeResult *stubWriteRequest
}

//nolint:ireturn // returning interface required by SDK interface
func (c *stubFGAWriteClient) Write(_ context.Context) gosdk.SdkClientWriteRequestInterface {
	return c.writeResult
}

func TestWriter_WriteTuple(t *testing.T) {
	t.Parallel()

	writer := NewWriter(&stubFGAWriteClient{writeResult: newStubWrite(nil)}, testLogger)

	err := writer.WriteTuple(context.Background(), "user:alice", "member", "system:global")
	if err != nil {
		t.Fatalf("WriteTuple() error = %v", err)
	}
}

func TestWriter_DeleteTuple(t *testing.T) {
	t.Parallel()

	writer := NewWriter(&stubFGAWriteClient{writeResult: newStubWrite(nil)}, testLogger)

	err := writer.DeleteTuple(context.Background(), "user:bob", "moderator", "system:global")
	if err != nil {
		t.Fatalf("DeleteTuple() error = %v", err)
	}
}

func TestWriter_WriteTuples(t *testing.T) {
	t.Parallel()

	writer := NewWriter(&stubFGAWriteClient{writeResult: newStubWrite(nil)}, testLogger)

	err := writer.WriteTuples(context.Background(), []TupleKey{
		{User: "user:alice", Relation: "author", Object: "post:1"},
		{User: "post:1", Relation: "global", Object: "system:global"},
	})
	if err != nil {
		t.Fatalf("WriteTuples() error = %v", err)
	}
}

func TestWriter_APIError(t *testing.T) {
	t.Parallel()

	writer := NewWriter(&stubFGAWriteClient{writeResult: newStubWrite(errFGA)}, testLogger)

	err := writer.WriteTuple(context.Background(), "user:alice", "member", "system:global")
	if err == nil {
		t.Fatal("WriteTuple() expected error, got nil")
	}
}
