package external

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/palantir/stacktrace"

	"go-poc/utils/activity"
	"go-poc/utils/log"
)

func NewMySQL(service string) (*sql.DB, error) {
	ctx := activity.NewContext("init_mysql")

	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbUsername := os.Getenv("MYSQL_USERNAME")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DB")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", dbUsername, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "can't open mysql connection"))
		return nil, stacktrace.Propagate(err, "can't open mysql connection")
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "can't ping mysql db"))
		return nil, stacktrace.Propagate(err, "can't ping mysql db")
	}

	initMigration(ctx, db, service, "mysql")

	return db, nil
}
