package config

// AuthenticatedAccount carries the minimal account information needed after a successful authentication.
type AuthenticatedAccount struct {
	ID   int64
	Role string
}
