package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"log"
	"os"
	"runtime/debug"
)

var snippetBoxDBName, snippetBoxDBUser, snippetboxDBPassword string
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

func init() {
	must(godotenv.Load())
	snippetBoxDBName = os.Getenv("DB_NAME")
	snippetBoxDBUser = os.Getenv("DB_USER")
	snippetboxDBPassword = os.Getenv("DB_PASSWORD")
}

func main() {
	rootCmd := buildRootCmd()
	must(rootCmd.Execute())
}

func buildRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Control the database lifecycle for the Snippetbox app",
		Run: func(cmd *cobra.Command, args []string) {
			db := getDB("postgres", "postgres", "postgres")
			Exec(db, fmt.Sprintf("DROP USER IF EXISTS %s", snippetBoxDBUser))
			Exec(db, fmt.Sprintf("DROP DATABASE IF EXISTS %s", snippetBoxDBName))
		},
	}
	cmd.AddCommand(buildMigrateUpCmd(), buildMigrateDownCmd())
	return cmd
}

func getDB(userName, password, dbName string) *sql.DB {
	dsn := fmt.Sprintf("host=localhost port=5432 user=%s password=%s sslmode=disable dbname=%s", userName, password, dbName)
	db, err := sql.Open("postgres", dsn)
	must(err)
	return db
}

func Exec(db *sql.DB, query string) {
	_, err := db.Exec(query)
	must(err)
}

func must(err error) {
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		errorLog.Output(2, trace)
		os.Exit(1)
	}
}
