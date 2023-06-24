package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/katerji/UserAuthKit/envs"
	"time"
)

type client struct {
	mysql *sql.DB
}

var instance *client

func GetDbInstance() *client {
	if instance == nil {
		instance, _ = getDbClient()
	}
	return instance
}

func getDbClient() (*client, error) {
	dbHost := envs.GetInstance().GetDbHost()
	dbUser := envs.GetInstance().GetDbUser()
	dbPort := envs.GetInstance().GetDbPort()
	dbPass := envs.GetInstance().GetDbPassword()
	dbName := envs.GetInstance().GetDbName()

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return &client{}, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return &client{
		db,
	}, nil
}

func (client *client) Fetch(query string, args ...any) []any {
	rows, err := client.mysql.Query(query, args...)
	if err != nil {
		fmt.Println(err.Error())
	}
	results := make([]interface{}, 0)
	for rows.Next() {
		var row interface{}

		err = rows.Scan(&row)
		results = append(results, row)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
	}
	return results
}

func (client *client) Insert(query string, args ...any) (int, error) {
	rows, err := client.mysql.Prepare(query)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	exec, err := rows.Exec(args...)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	insertId, err := exec.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	return int(insertId), nil
}

func (client *client) Exec(query string, args ...any) bool {
	rows, err := client.mysql.Prepare(query)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = rows.Exec(args...)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func (client *client) QueryRow(query string, args ...any) *sql.Row {
	return client.mysql.QueryRow(query, args...)
}

func (client *client) Query(query string, args ...any) (*sql.Rows, error) {
	return client.mysql.Query(query, args...)
}

func (client *client) Ping() error {
	return client.mysql.Ping()
}
