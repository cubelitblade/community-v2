// package main starts the authz API server for the community-v2 backend
package main

import (
	"log"

	"github.com/cubelitblade/community-v2/backend/services/authz/internal/bootstrap"
)

func main() {
	app := bootstrap.NewApp()
	app.Run()

	log.Println("server stopped")
}
