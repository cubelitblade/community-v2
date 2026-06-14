// package main starts the community-v2 service
package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/cubelitblade/community-v2/backend/services/post/internal/bootstrap"
	"github.com/cubelitblade/community-v2/backend/services/post/migrations"
)

func main() {
	migrateOnly := flag.Bool("migrate", false, "Run database migrations.")

	flag.Parse()

	if *migrateOnly {
		uri := os.Getenv("POST_MIGRATE_URI")

		db, err := sql.Open("pgx", uri)
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

	bootstrap.NewApp().Run()
}
