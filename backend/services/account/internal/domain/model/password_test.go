package model_test

import (
	"errors"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/cubelitblade/community-v2/backend/services/account/internal/domain/model"
)

const (
	validPwd = "this-is-a-valid-password"
	emptyPwd = ""
	padChar  = "a"
)

//nolint:gochecknoglobals // test data requiring runtime initialization
var tooLongPassword = strings.Repeat(padChar, model.MaxPasswordBytes+1)

func TestValidatePassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		wantErr error
		name    string
		raw     string
	}{
		{
			name:    "rejects empty password",
			raw:     emptyPwd,
			wantErr: model.ErrPasswordEmpty,
		},
		{
			name:    "rejects password shorter than minimum",
			raw:     padChar,
			wantErr: model.ErrPasswordTooShort,
		},
		{
			name:    "accepts valid password",
			raw:     validPwd,
			wantErr: nil,
		},
		{
			name:    "rejects password longer than bcrypt byte limit",
			raw:     tooLongPassword,
			wantErr: model.ErrPasswordTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotErr := model.ValidatePassword(tt.raw)
			assertPasswordError(t, gotErr, tt.wantErr)
		})
	}
}

func TestHashPassword(t *testing.T) {
	t.Parallel()

	raw := validPwd

	hash, err := model.HashPassword(raw)
	if err != nil {
		t.Fatalf("HashPassword() failed: %v", err)
	}

	if hash.Value() == emptyPwd {
		t.Fatal("HashPassword() returned empty hash")
	}

	if hash.Value() == raw {
		t.Fatal("HashPassword() returned raw password")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(hash.Value()), []byte(raw),
	); err != nil {
		t.Fatalf("generated hash does not match raw password: %v", err)
	}
}

func TestPasswordHashFromStorage(t *testing.T) {
	t.Parallel()

	t.Run("accepts valid hash string", func(t *testing.T) {
		t.Parallel()

		hash, err := model.PasswordHashFromStorage("some-hash")
		if err != nil {
			t.Fatalf("PasswordHashFromStorage() failed: %v", err)
		}

		if hash.Value() != "some-hash" {
			t.Errorf("Value() = %q, want %q", hash.Value(), "some-hash")
		}
	})

	t.Run("rejects empty hash string", func(t *testing.T) {
		t.Parallel()

		_, err := model.PasswordHashFromStorage("")
		if !errors.Is(err, model.ErrPasswordHashEmpty) {
			t.Errorf("error = %v, want %v", err, model.ErrPasswordHashEmpty)
		}
	})
}

func TestPasswordHashValue(t *testing.T) {
	t.Parallel()

	hash, err := model.HashPassword(validPwd)
	if err != nil {
		t.Fatalf("HashPassword() failed: %v", err)
	}

	if hash.Value() == emptyPwd {
		t.Fatal("Value() returned empty string")
	}
}

func TestPasswordHashString(t *testing.T) {
	t.Parallel()

	hash, err := model.HashPassword(validPwd)
	if err != nil {
		t.Fatalf("HashPassword() failed: %v", err)
	}

	if hash.String() != "[redacted]" {
		t.Errorf("String() = %q, want %q", hash.String(), "[redacted]")
	}
}

func assertPasswordError(t *testing.T, gotErr, wantErr error) {
	t.Helper()

	if wantErr == nil {
		if gotErr != nil {
			t.Errorf("unexpected error: %v", gotErr)
		}

		return
	}

	if gotErr == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(gotErr, wantErr) {
		t.Errorf("error = %v, want %v", gotErr, wantErr)
	}
}
