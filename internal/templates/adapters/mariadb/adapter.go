package mariadb

import (
	"database/sql"
	"fmt"
)

type MariaDB struct {
	DB *sql.DB
}

func New(host, port, database, username, password string) (*MariaDB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, database)
	client, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	// Ping to see if connection was successful
	err = client.Ping()
	if err != nil {
		return nil, err
	}

	return &MariaDB{
		DB: client,
	}, nil
}

func (m *MariaDB) Close() {
	if m.DB != nil {
		m.DB.Close()
	}
}
