package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func buildMigrateUpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Initializes the database and the database user used by the Snippetbox app",
		Run: func(cmd *cobra.Command, args []string) {
			db := getDB("postgres", "postgres", "postgres")
			Exec(db, fmt.Sprintf("CREATE DATABASE %s", snippetBoxDBName))
			infoLog.Printf("Created Postgres database %s", snippetBoxDBName)
			must(db.Close())

			db = getDB("postgres", "postgres", snippetBoxDBName)
			script, err := os.ReadFile("./cmd/migration/setup.sql")
			must(err)
			infoLog.Printf("Executed setup.sql")
			Exec(db, string(script))
			Exec(db, fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", snippetBoxDBUser, snippetboxDBPassword))
			Exec(db, fmt.Sprintf("GRANT SELECT, INSERT, UPDATE, DELETE ON snippets, sessions, users TO %s", snippetBoxDBUser))
			infoLog.Printf("Create Postgres user %s", snippetBoxDBUser)
			must(db.Close())

		},
	}
	return cmd
}
