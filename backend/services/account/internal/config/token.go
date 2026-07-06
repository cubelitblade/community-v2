package config

import (
	"time"
)

// TokenConfig holds the configuration for token issuance and validity periods.
type TokenConfig struct {
	TokenIssuer         string        `koanf:"issuer"`
	AccessTokenValidity time.Duration `koanf:"access_token_validity"`
}
