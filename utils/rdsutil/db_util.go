package rdsutil

import (
	"database/sql"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type DBClient struct {
	Db *sql.DB
}

func NewDBClient(dbHost, dbPort, dbUser, dbPassword string) (*DBClient, error) {
	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/")
	if err != nil {
		return nil, err
	}

	return &DBClient{Db: db}, nil
}

func (dbc *DBClient) ExecuteSQL(sqlFilePath string) error {
	queryBytes, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return err
	}
	queries := strings.Split(string(queryBytes), ";")
	for _, query := range queries {
		trimmedQuery := strings.TrimSpace(query)
		if trimmedQuery != "" {
			_, err = dbc.Db.Exec(trimmedQuery)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
