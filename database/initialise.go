package database

import (
	"database/sql"
	"fmt"
	"untitled-agnostic-sql-migration-package/structs"
)

func InitialiseDatabase(p *structs.ProjectConfig) {
	psqlDsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.Username, p.Password, p.Database)

	db, err := sql.Open("postgres", psqlDsn)
	if err != nil {
		fmt.Println("Failed to connect to database for project: " + p.Name + " due to: " + err.Error())
		return
	}

	defer db.Close()

	sqlStatement := `
		create table "_sqltoolHistory"
		(
			id serial,
			filename text,
			migrated timestamp default now(),
			batch_number int default 0 not null
		);
		
		create unique index _sqltoolhistory_id_uindex
			on "_sqltoolHistory" (id);
		
		alter table "_sqltoolHistory"
			add constraint _sqltoolhistory_pk
				primary key (id);
	`

	_, err = db.Exec(sqlStatement)

	if err != nil {
		panic(err)
	}
}
