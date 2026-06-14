package bootstrap

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/cubelitblade/community-v2/backend/services/post/internal/shared"
)

func LoadConfig() (shared.Config, error) {
	k := koanf.New(".")

	// Load from TOML.
	if tomlPath := os.Getenv("POST_CONFIG_TOML_PATH"); tomlPath != "" {
		if err := k.Load(file.Provider(tomlPath), toml.Parser()); err != nil {
			return shared.Config{}, fmt.Errorf("load config from toml: %w", err)
		}
	}

	// Load from environment variables, will override duplicated config.
	if err := k.Load(env.Provider(".", env.Opt{
		Prefix: "POST_",
		TransformFunc: func(k, v string) (string, any) {
			k = strings.ReplaceAll(strings.ToLower(k), "_", ".")
			if strings.Contains(v, " ") {
				return k, strings.Split(v, " ")
			}

			return k, v
		},
	}), nil); err != nil {
		return shared.Config{}, fmt.Errorf("load config from environment variables: %w", err)
	}

	var config shared.Config
	k.Unmarshal("post", &config)

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(config); err != nil {
		return shared.Config{}, fmt.Errorf("validate config: %w", err)
	}

	return config, nil
}
