package events_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	pkgEvents "github.com/cubelitblade/community-v2/backend/pkg/events"
	"github.com/cubelitblade/community-v2/backend/services/authz/internal/events"
)

const (
	keyUserID    = "user_id"
	keyRole      = "role"
	keyAccountID = "account_id"
	keyPostID    = "post_id"
)

const (
	testUserID    = "42"
	testAccountID = "acct1"
	testPostID    = "post1"
	testRole      = "moderator"
)

type deleteCall struct {
	User     string
	Relation string
	Object   string
}

type stubTupleWriter struct {
	tuples     []events.TupleKey
	deleteCall *deleteCall
}

func (writer *stubTupleWriter) WriteTuple(_ context.Context, user, relation, object string) error {
	writer.tuples = append(writer.tuples, events.TupleKey{User: user, Relation: relation, Object: object})

	return nil
}

func (writer *stubTupleWriter) WriteTuples(_ context.Context, tuples []events.TupleKey) error {
	writer.tuples = append(writer.tuples, tuples...)

	return nil
}

func (writer *stubTupleWriter) DeleteTuple(_ context.Context, user, relation, object string) error {
	writer.deleteCall = &deleteCall{User: user, Relation: relation, Object: object}

	return nil
}

func mustMarshal(t *testing.T, v any) []byte {
	t.Helper()

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	return data
}

func assertTuple(t *testing.T, tuple events.TupleKey, user, relation, object string) {
	t.Helper()

	if tuple.User != user {
		t.Errorf("User = %q, want %q", tuple.User, user)
	}

	if tuple.Relation != relation {
		t.Errorf("Relation = %q, want %q", tuple.Relation, relation)
	}

	if tuple.Object != object {
		t.Errorf("Object = %q, want %q", tuple.Object, object)
	}
}

func newRegisteredPayload() map[string]string {
	return map[string]string{keyUserID: testUserID}
}

func newRoleAssignedPayload() map[string]string {
	return map[string]string{keyUserID: testUserID, keyRole: testRole}
}

func newAccountCreatedPayload() map[string]string {
	return map[string]string{keyUserID: testUserID, keyAccountID: testAccountID}
}

func newPostCreatedPayload() map[string]string {
	return map[string]string{keyUserID: testUserID, keyPostID: testPostID}
}

func newStubWriter() *stubTupleWriter {
	return &stubTupleWriter{
		tuples:     nil,
		deleteCall: nil,
	}
}

func TestHandler_UserRegistered(t *testing.T) {
	t.Parallel()

	writer := newStubWriter()
	handler := events.NewHandler(writer)

	body := mustMarshal(t, newRegisteredPayload())

	err := handler.Handle(context.Background(), "user.registered", body)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	if len(writer.tuples) != 1 {
		t.Fatalf("len(tuples) = %d, want 1", len(writer.tuples))
	}

	assertTuple(t, writer.tuples[0], "user:"+testUserID, "member", "system:global")
}

func TestHandler_UserRoleAssigned(t *testing.T) {
	t.Parallel()

	writer := newStubWriter()
	handler := events.NewHandler(writer)

	body := mustMarshal(t, newRoleAssignedPayload())

	err := handler.Handle(context.Background(), "user.role_assigned", body)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	if len(writer.tuples) != 1 {
		t.Fatalf("len(tuples) = %d, want 1", len(writer.tuples))
	}

	assertTuple(t, writer.tuples[0], "user:"+testUserID, testRole, "system:global")
}

func TestHandler_UserRoleRevoked(t *testing.T) {
	t.Parallel()

	writer := newStubWriter()
	handler := events.NewHandler(writer)

	body := mustMarshal(t, newRoleAssignedPayload())

	err := handler.Handle(context.Background(), "user.role_revoked", body)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	if writer.deleteCall == nil {
		t.Fatal("expected DeleteTuple call")
	}

	if writer.deleteCall.User != "user:"+testUserID {
		t.Errorf("User = %q, want %q", writer.deleteCall.User, "user:"+testUserID)
	}

	if writer.deleteCall.Relation != testRole {
		t.Errorf("Relation = %q, want %q", writer.deleteCall.Relation, testRole)
	}

	if writer.deleteCall.Object != "system:global" {
		t.Errorf("Object = %q, want %q", writer.deleteCall.Object, "system:global")
	}
}

func TestHandler_AccountCreated(t *testing.T) {
	t.Parallel()

	writer := newStubWriter()
	handler := events.NewHandler(writer)

	body := mustMarshal(t, newAccountCreatedPayload())

	err := handler.Handle(context.Background(), "account.created", body)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	if len(writer.tuples) != 2 {
		t.Fatalf("len(tuples) = %d, want 2", len(writer.tuples))
	}

	assertTuple(t, writer.tuples[0], "user:"+testUserID, "owner", "account:"+testAccountID)
	assertTuple(t, writer.tuples[1], "account:"+testAccountID, "global", "system:global")
}

func TestHandler_PostCreated(t *testing.T) {
	t.Parallel()

	writer := newStubWriter()
	handler := events.NewHandler(writer)

	body := mustMarshal(t, newPostCreatedPayload())

	err := handler.Handle(context.Background(), "post.created", body)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	if len(writer.tuples) != 2 {
		t.Fatalf("len(tuples) = %d, want 2", len(writer.tuples))
	}

	assertTuple(t, writer.tuples[0], "user:"+testUserID, "author", "post:"+testPostID)
	assertTuple(t, writer.tuples[1], "post:"+testPostID, "global", "system:global")
}

func TestHandler_UnknownEvent(t *testing.T) {
	t.Parallel()

	writer := newStubWriter()
	handler := events.NewHandler(writer)

	err := handler.Handle(context.Background(), "unknown.event", []byte(`{}`))
	if !errors.Is(err, pkgEvents.ErrUnknownEvent) {
		t.Errorf("error = %v, want %v", err, pkgEvents.ErrUnknownEvent)
	}
}

func TestHandler_InvalidPayload(t *testing.T) {
	t.Parallel()

	writer := newStubWriter()
	handler := events.NewHandler(writer)

	err := handler.Handle(context.Background(), "user.registered", []byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid payload")
	}
}
