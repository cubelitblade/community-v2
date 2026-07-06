package config

// CookieConfig holds the configuration for HTTP cookies.
type CookieConfig struct {
	Secure   bool   `koanf:"secure"`
	SameSite string `koanf:"same_site"`
}
