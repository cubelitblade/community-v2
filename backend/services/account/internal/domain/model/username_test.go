package model_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/cubelitblade/community-v2/backend/services/account/internal/domain/model"
)

const (
	fill              string     = "a"
	testID            model.ID = 1
	maxArchiveID      model.ID = 9223372036854775807
	maxArchivedSuffix            = "#archived_9223372036854775807"
)

type usernameCase struct {
	wantErr error
	name    string
	value   string
	want    model.Username
}

type archiveUsernameCase struct {
	wantErr  error
	name     string
	username model.Username
	want     model.Username
	id       model.ID
}

func TestNewUsername(t *testing.T) {
	t.Parallel()

	runUsernameCases(t, "NewUsername", model.NewUsername, newUsernameCases())
}

func TestUsernameFromStorage(t *testing.T) {
	t.Parallel()

	runUsernameCases(
		t,
		"UsernameFromStorage",
		model.UsernameFromStorage,
		usernameFromStorageCases(),
	)
}

func TestArchiveUsername(t *testing.T) {
	t.Parallel()

	for _, testCase := range archiveUsernameCases() {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := model.ArchiveUsername(
				testCase.username,
				testCase.id,
			)
			assertUsernameResult(
				t,
				"ArchiveUsername",
				got,
				gotErr,
				testCase.want,
				testCase.wantErr,
			)
		})
	}
}

func runUsernameCases(
	t *testing.T,
	funcName string,
	factory func(string) (model.Username, error),
	tests []usernameCase,
) {
	t.Helper()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := factory(testCase.value)
			assertUsernameResult(
				t,
				funcName,
				got,
				gotErr,
				testCase.want,
				testCase.wantErr,
			)
		})
	}
}

func assertUsernameResult(
	t *testing.T,
	funcName string,
	got model.Username,
	gotErr error,
	want model.Username,
	wantErr error,
) {
	t.Helper()

	if wantErr != nil {
		assertUsernameError(t, funcName, gotErr, wantErr)

		return
	}

	if gotErr != nil {
		t.Fatalf("%s() failed: %v", funcName, gotErr)
	}

	if got != want {
		t.Errorf("%s() = %v, want %v", funcName, got, want)
	}
}

func assertUsernameError(t *testing.T, funcName string, gotErr, wantErr error) {
	t.Helper()

	if gotErr == nil {
		t.Fatalf("%s() succeeded unexpectedly", funcName)
	}

	if !errors.Is(gotErr, wantErr) {
		t.Errorf("%s() error = %v, want %v", funcName, gotErr, wantErr)
	}
}

func newUsernameCases() []usernameCase {
	minUsername := strings.Repeat(fill, model.MinUsernameLength)
	maxUsername := strings.Repeat(fill, model.MaxUsernameLength)

	return []usernameCase{
		{
			name:    "accepts minimum length username",
			value:   minUsername,
			want:    model.Username(minUsername),
			wantErr: nil,
		},
		{
			name:    "accepts maximum length username",
			value:   maxUsername,
			want:    model.Username(maxUsername),
			wantErr: nil,
		},
		{
			name:    "rejects username shorter than minimum length",
			value:   fill,
			want:    "",
			wantErr: model.ErrUsernameLength,
		},
		{
			name:    "rejects username longer than maximum length",
			value:   strings.Repeat(fill, model.MaxUsernameLength+1),
			want:    "",
			wantErr: model.ErrUsernameLength,
		},
		{
			name:    "trims surrounding spaces",
			value:   " Alice ",
			want:    "Alice",
			wantErr: nil,
		},
		{
			name:    "rejects blank username",
			value:   "   ",
			want:    "",
			wantErr: model.ErrUsernameBlank,
		},
	}
}

func usernameFromStorageCases() []usernameCase {
	tooLongNormalUsername := strings.Repeat(fill, model.MaxUsernameLength+1)
	tooLongStoredUsername := strings.Repeat(
		fill,
		model.MaxStoredUsernameLength+1,
	)

	return []usernameCase{
		{
			name:    "accepts normal username",
			value:   "bob",
			want:    "bob",
			wantErr: nil,
		},
		{
			name:    "accepts archived username",
			value:   "Alice#archived_123",
			want:    "Alice#archived_123",
			wantErr: nil,
		},
		{
			name:    "accepts username longer than normal limit",
			value:   tooLongNormalUsername,
			want:    model.Username(tooLongNormalUsername),
			wantErr: nil,
		},
		{
			name:    "rejects blank stored username",
			value:   "   ",
			want:    "",
			wantErr: model.ErrUsernameBlank,
		},
		{
			name:    "rejects username longer than storage limit",
			value:   tooLongStoredUsername,
			want:    "",
			wantErr: model.ErrUsernameLength,
		},
	}
}

func archiveUsernameCases() []archiveUsernameCase {
	maxUsername := strings.Repeat(fill, model.MaxUsernameLength)
	tooLongStoredUsername := strings.Repeat(
		fill,
		model.MaxStoredUsernameLength,
	)

	return []archiveUsernameCase{
		{
			name:     "generates archived username from normal username",
			username: "Alice",
			id:       maxArchiveID,
			want:     "Alice" + maxArchivedSuffix,
			wantErr:  nil,
		},
		{
			name:     "keeps archived username within storage limit",
			username: model.Username(maxUsername),
			id:       maxArchiveID,
			want:     model.Username(maxUsername + maxArchivedSuffix),
			wantErr:  nil,
		},
		{
			name:     "rejects archived username longer than storage limit",
			username: model.Username(tooLongStoredUsername),
			id:       testID,
			want:     "",
			wantErr:  model.ErrUsernameLength,
		},
	}
}
