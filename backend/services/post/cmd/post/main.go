// package main starts the community-v2 service
package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/cubelitblade/community-v2/backend/services/post/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	migrateOnly := flag.Bool("migrate", false, "Run database migrations.")
	flag.Parse()

	if *migrateOnly {
		dsn := os.Getenv("POST_MIGRATE_DSN")
		db, err := sql.Open("pgx", dsn)
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}
		defer db.Close()

		log.Println("Running database migrations...")
		if err := migrations.Up(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migrations completed successfully.")
		return
	}

	log.Println("service has not been implemented yet")
}
