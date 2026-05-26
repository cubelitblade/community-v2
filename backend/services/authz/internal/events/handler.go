// Package events provides authz-specific event handling: routing keys, payload types,
// and OpenFGA tuple mapping. Generic AMQP infrastructure lives in pkg/events.
package events

import (
	"context"
	"encoding/json"
	"fmt"

	pkgEvents "github.com/cubelitblade/community-v2/backend/pkg/events"
)

// Payload types for domain events.

//nolint:tagliatelle // JSON keys are cross-service contract, not Go conventions
type userRegisteredPayload struct {
	UserID string `json:"user_id"`
}

//nolint:tagliatelle // JSON keys are cross-service contract, not Go conventions
type userRoleAssignedPayload struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

//nolint:tagliatelle // JSON keys are cross-service contract, not Go conventions
type userRoleRevokedPayload struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

//nolint:tagliatelle // JSON keys are cross-service contract, not Go conventions
type accountCreatedPayload struct {
	UserID    string `json:"user_id"`
	AccountID string `json:"account_id"`
}

//nolint:tagliatelle // JSON keys are cross-service contract, not Go conventions
type postCreatedPayload struct {
	UserID string `json:"user_id"`
	PostID string `json:"post_id"`
}

// routing keys for domain events.
const (
	rkUserRegistered   = "user.registered"
	rkUserRoleAssigned = "user.role_assigned"
	rkUserRoleRevoked  = "user.role_revoked"
	rkAccountCreated   = "account.created"
	rkPostCreated      = "post.created"
)

type handlerFunc func(ctx context.Context, body []byte) error

// TupleWriter is the interface for writing relationship tuples to the authorization store.
type TupleWriter interface {
	WriteTuple(ctx context.Context, user, relation, object string) error
	DeleteTuple(ctx context.Context, user, relation, object string) error
	WriteTuples(ctx context.Context, tuples []TupleKey) error
}

// Handler dispatches incoming events to the appropriate tuple writer operations.
type Handler struct {
	writer TupleWriter
	routes map[string]handlerFunc
}

// NewHandler creates an event Handler backed by the given TupleWriter.
func NewHandler(writer TupleWriter) *Handler {
	if writer == nil {
		panic("nil writer")
	}

	handler := &Handler{
		writer: writer,
		routes: make(map[string]handlerFunc),
	}

	handler.routes[rkUserRegistered] = handler.handleUserRegistered
	handler.routes[rkUserRoleAssigned] = handler.handleUserRoleAssigned
	handler.routes[rkUserRoleRevoked] = handler.handleUserRoleRevoked
	handler.routes[rkAccountCreated] = handler.handleAccountCreated
	handler.routes[rkPostCreated] = handler.handlePostCreated

	return handler
}

// Handle dispatches a raw event to the appropriate handler based on the routing key.
func (h *Handler) Handle(ctx context.Context, routingKey string, body []byte) error {
	fn, ok := h.routes[routingKey]
	if !ok {
		return fmt.Errorf("%w: %s", pkgEvents.ErrUnknownEvent, routingKey)
	}

	if err := fn(ctx, body); err != nil {
		return fmt.Errorf("handle event %s: %w", routingKey, err)
	}

	return nil
}

func (h *Handler) handleUserRegistered(ctx context.Context, body []byte) error {
	var payload userRegisteredPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("unmarshal user.registered: %w", err)
	}

	if err := h.writer.WriteTuple(ctx, userKey(payload.UserID), "member", "system:global"); err != nil {
		return fmt.Errorf("write tuple: %w", err)
	}

	return nil
}

func (h *Handler) handleUserRoleAssigned(ctx context.Context, body []byte) error {
	var payload userRoleAssignedPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("unmarshal user.role_assigned: %w", err)
	}

	if err := h.writer.WriteTuple(ctx, userKey(payload.UserID), payload.Role, "system:global"); err != nil {
		return fmt.Errorf("write tuple: %w", err)
	}

	return nil
}

func (h *Handler) handleUserRoleRevoked(ctx context.Context, body []byte) error {
	var payload userRoleRevokedPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("unmarshal user.role_revoked: %w", err)
	}

	if err := h.writer.DeleteTuple(ctx, userKey(payload.UserID), payload.Role, "system:global"); err != nil {
		return fmt.Errorf("delete tuple: %w", err)
	}

	return nil
}

func (h *Handler) handleAccountCreated(ctx context.Context, body []byte) error {
	var payload accountCreatedPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("unmarshal account.created: %w", err)
	}

	if err := h.writer.WriteTuples(ctx, []TupleKey{
		{User: userKey(payload.UserID), Relation: "owner", Object: accountKey(payload.AccountID)},
		{User: accountKey(payload.AccountID), Relation: "global", Object: "system:global"},
	}); err != nil {
		return fmt.Errorf("write tuples: %w", err)
	}

	return nil
}

func (h *Handler) handlePostCreated(ctx context.Context, body []byte) error {
	var payload postCreatedPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("unmarshal post.created: %w", err)
	}

	if err := h.writer.WriteTuples(ctx, []TupleKey{
		{User: userKey(payload.UserID), Relation: "author", Object: postKey(payload.PostID)},
		{User: postKey(payload.PostID), Relation: "global", Object: "system:global"},
	}); err != nil {
		return fmt.Errorf("write tuples: %w", err)
	}

	return nil
}

func userKey(id string) string    { return "user:" + id }
func accountKey(id string) string { return "account:" + id }
func postKey(id string) string    { return "post:" + id }
