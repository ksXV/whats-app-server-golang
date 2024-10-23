package database

import (
	"database/sql"
	"log"

	_ "github.com/marcboeker/go-duckdb"
)

type DBConnection struct {
	DB *sql.DB
}

const DATABASE_PATH = "database/server.db"

func NewDBConnection() DBConnection {
	dbConn := DBConnection{}

	dbConn.initDB()
	dbConn.createTables()

	return dbConn
}

func (dbCon *DBConnection) initDB() {
	db, err := sql.Open("duckdb", DATABASE_PATH)
	if err != nil {
		log.Fatal(err.Error())
	}

	dbCon.DB = db
}

func (dbConn *DBConnection) createTables() {
	if dbConn.DB == nil {
		log.Fatal("DATABASE IS NULL!! ABORTING")
	}

	userTable := `
		CREATE TABLE IF NOT EXISTS user (
			id UUID PRIMARY KEY default uuid(),
			username varchar NOT NULL,
			profilePicture varchar NOT NULL,
			createdAt date NOT NULL,
			usesOAuth boolean NOT NULL
		);
	`
	userCredentialsTable := `
		CREATE TABLE IF NOT EXISTS user_credentials (
			id UUID PRIMARY KEY DEFAULT uuid(),
			user_id UUID NOT NULL,
			email varchar UNIQUE NOT NULL,
			hash varchar UNIQUE NOT NULL,
			salt varchar NOT NULL,
			FOREIGN KEY (user_id) REFERENCES user(id)
		);
	`
	userOAuthTable := `
		CREATE TABLE IF NOT EXISTS user_oauth (
			id uuid PRIMARY KEY DEFAULT uuid(),
			user_id uuid NOT NULL,
			email varchar NOT NULL,
			provider varchar NOT NULL,
			FOREIGN KEY (user_id) REFERENCES user(id),
			UNIQUE(email, provider)
		);
	`

	_, err := dbConn.DB.Exec(userTable + userOAuthTable + userCredentialsTable)
	if err != nil {
		log.Fatal(err.Error())
	}
}
