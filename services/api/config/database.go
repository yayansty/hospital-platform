package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/godror/godror"
	_ "github.com/microsoft/go-mssqldb"
)

// ConnectDB connects to SQL Server.
func ConnectDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	connString := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable",
		user,
		password,
		host,
		port,
		database,
	)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// ConnectOracle connects to Oracle.
func ConnectOracle() (*sql.DB, error) {
	host := os.Getenv("ORACLE_HOST")
	port := os.Getenv("ORACLE_PORT")
	service := os.Getenv("ORACLE_SERVICE")
	user := os.Getenv("ORACLE_USER")
	password := os.Getenv("ORACLE_PASSWORD")

	connectString := fmt.Sprintf(
		"%s:%s/%s",
		host,
		port,
		service,
	)

	dsn := fmt.Sprintf(
		`user="%s" password="%s" connectString="%s"`,
		user,
		password,
		connectString,
	)

	db, err := sql.Open("godror", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
