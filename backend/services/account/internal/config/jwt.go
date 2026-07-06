package config

import (
	jwt2 "github.com/golang-jwt/jwt/v5"
)

// JWTConfig holds the configuration for JSON Web Token authentication.
type JWTConfig struct {
	Secret string `koanf:"secret"`
}

// AccessTokenClaims defines the JWT claims for an access token, including the user's role.
type AccessTokenClaims struct {
	jwt2.RegisteredClaims

	Role string `json:"role"`
}
