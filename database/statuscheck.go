package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"untitled-agnostic-sql-migration-package/structs"
)

/*
	Initialise a connection to the given database config and check whether it has a
	table called _sqltoolHistory, if it does we'll consider it initialised, if not (or if the db isn't
	reachable for whatever reason), then we'll consider it not initialised
 */
func CheckInitialisationStatus(p *structs.ProjectConfig, c chan structs.DatabaseState) {
	var state structs.DatabaseState
	if p.Driver != "postgres" {
		fmt.Println("Unsupported database type")
		c <- state
		return
	}

	psqlDsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.Username, p.Password, p.Database)

	db, err := sql.Open("postgres", psqlDsn)
	if err != nil {
		fmt.Println("Failed to connect to database for project: " + p.Name + " due to: " + err.Error())
		c <- state
		return
	}

	defer db.Close()
	state.Reachable = true

	sqlStatement := `
	SELECT EXISTS (
		SELECT FROM pg_tables
		WHERE schemaname = 'public'
		AND tablename = '_sqltoolHistory'
	)`

	err = db.QueryRow(sqlStatement).Scan(&state.Initialised)

	if err != nil {
		fmt.Println(err)
		c <- state
	}

	c <- state
}
