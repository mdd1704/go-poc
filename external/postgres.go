package external

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/palantir/stacktrace"

	"go-poc/utils/activity"
	"go-poc/utils/log"
)

func NewPostgres(service string) (*sql.DB, error) {
	ctx := activity.NewContext("init_postgres")

	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUsername := os.Getenv("POSTGRES_USERNAME")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUsername, dbPassword, dbName)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "can't open postgres connection"))
		return nil, stacktrace.Propagate(err, "can't open postgres connection")
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "can't ping postgres db"))
		return nil, stacktrace.Propagate(err, "can't ping postgres db")
	}

	initMigration(ctx, db, service, "postgres")

	return db, nil
}
