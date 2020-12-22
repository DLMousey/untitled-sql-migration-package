package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"untitled-agnostic-sql-migration-package/structs"
)

func Migrate(p *structs.ProjectConfig) {
	psqlDsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.Username, p.Password, p.Database)

	db, err := sql.Open("postgres", psqlDsn)
	if err != nil {
		fmt.Println("Failed to connect to database for project: " + p.Name + " due to: " + err.Error())
		return
	}

	defer db.Close()

	root := p.Directory + "/" + p.SqlDirectory
	files, err := ioutil.ReadDir(root)
	if err != nil {
		panic(err)
	}

	fmt.Println("Processing migrations for project: " + p.Name)
	var sqlStatement string
	for _, file := range files {
		fmt.Println("Located migration: " + file.Name())

		var count int
		sqlStatement = "SELECT COUNT(*) FROM \"_sqltoolHistory\" WHERE filename = $1"
		err := db.QueryRow(sqlStatement, file.Name()).Scan(&count)

		if err != nil {
			panic(err)
		}

		if count > 0 {
			fmt.Println("Migration " + file.Name() + " already run, skipping")
		} else {
			sqlStatement, err := ioutil.ReadFile(root + "/" + file.Name())
			if err != nil {
				panic(err)
			}

			fmt.Println("Executing migration " + file.Name() + " for project " + p.Name)
			_, err = db.Exec(string(sqlStatement))

			batchStatement := "INSERT INTO \"_sqltoolHistory\" VALUES(default, $1, now(), 1)"
			_, err = db.Exec(batchStatement, file.Name())

			if err != nil {
				panic(err)
			}
		}
	}

}