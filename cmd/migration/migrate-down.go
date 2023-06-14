package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func buildMigrateDownCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down",
		Short: "Teardown the database and the database user used by the Snippetbox app",
		Run: func(cmd *cobra.Command, args []string) {
			db := getDB("postgres", "postgres", "postgres")
			Exec(db, fmt.Sprintf("DROP DATABASE IF EXISTS %s", snippetBoxDBName))
			infoLog.Printf("Dropped Postgres database %s", snippetBoxDBName)
			Exec(db, fmt.Sprintf("DROP USER IF EXISTS %s", snippetBoxDBUser))
			infoLog.Printf("Dropped Postgres user %s", snippetBoxDBUser)
		},
	}
	return cmd
}
