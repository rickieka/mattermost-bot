package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"git.rickiekarp.net/rickie/mattermost-bot/data/config"
)

var db *sql.DB

func GetResultSet(sqlConfig config.SqlConfig, query string) (*sql.Rows, error) {

	// Open up our database connection.
	db, err := sql.Open("mysql", sqlConfig.User+":"+sqlConfig.Password+"@tcp("+sqlConfig.Host+")/"+sqlConfig.Database)

	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	// Execute the query
	results, err := db.Query(query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// Return the result set
	return results, err
}
