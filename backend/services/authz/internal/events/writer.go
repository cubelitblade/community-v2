package events

import (
	"context"
	"fmt"
	"log/slog"

	gosdk "github.com/openfga/go-sdk/client"
)

// fgaWriteClient is the subset of the OpenFGA Go SDK client that Writer needs.
type fgaWriteClient interface {
	Write(ctx context.Context) gosdk.SdkClientWriteRequestInterface
}

// TupleKey represents a relationship tuple to write to OpenFGA.
type TupleKey struct {
	User     string
	Relation string
	Object   string
}

// Writer writes relationship tuples to OpenFGA.
type Writer struct {
	client fgaWriteClient
	logger *slog.Logger
}

// NewWriter creates a Writer backed by the given OpenFGA SDK client.
func NewWriter(client fgaWriteClient, logger *slog.Logger) *Writer {
	if client == nil {
		panic("nil client")
	}

	if logger == nil {
		panic("nil logger")
	}

	return &Writer{
		client: client,
		logger: logger,
	}
}

// WriteTuple writes a single relationship tuple to OpenFGA.
func (w *Writer) WriteTuple(ctx context.Context, user, relation, object string) error {
	return w.WriteTuples(ctx, []TupleKey{{User: user, Relation: relation, Object: object}})
}

// DeleteTuple removes a single relationship tuple from OpenFGA.
func (w *Writer) DeleteTuple(ctx context.Context, user, relation, object string) error {
	_, err := w.client.Write(ctx).Body(gosdk.ClientWriteRequest{
		Deletes: []gosdk.ClientTupleKeyWithoutCondition{
			{User: user, Relation: relation, Object: object},
		},
	}).Execute()
	if err != nil {
		return fmt.Errorf("openfga delete: %w", err)
	}

	return nil
}

// WriteTuples writes multiple relationship tuples to OpenFGA in a single request.
func (w *Writer) WriteTuples(ctx context.Context, tuples []TupleKey) error {
	keys := make([]gosdk.ClientTupleKey, len(tuples))
	for i, t := range tuples {
		keys[i] = gosdk.ClientTupleKey{
			User:     t.User,
			Relation: t.Relation,
			Object:   t.Object,
		}
	}

	_, err := w.client.Write(ctx).Body(gosdk.ClientWriteRequest{
		Writes: keys,
	}).Execute()
	if err != nil {
		return fmt.Errorf("openfga write: %w", err)
	}

	return nil
}

// Ensure *OpenFgaClient satisfies fgaWriteClient.
var _ fgaWriteClient = (*gosdk.OpenFgaClient)(nil)
