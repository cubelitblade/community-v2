package config

import (
	"time"
)

// GormConfig holds the configuration for the GORM.
type GormConfig struct {
	SkipDefaultTransaction                   bool              `koanf:"skip_default_transaction"`
	PrepareStmt                              bool              `koanf:"prepare_stmt"`
	DisableForeignKeyConstraintWhenMigrating bool              `koanf:"disable_foreign_key_constraint_when_migrating"`
	Logger                                   *GormLoggerConfig `koanf:"logger"`
}

// GormLoggerConfig holds the configuration for the GORM logger.
type GormLoggerConfig struct {
	Level                     string        `koanf:"level"`
	SlowThreshold             time.Duration `koanf:"slow_threshold"`
	IgnoreRecordNotFoundError bool          `koanf:"ignore_record_not_found_error"`
	ParameterizedQuery        bool          `koanf:"parameterized_query"`
	Colorful                  bool          `koanf:"colorful"`
}
