package external

import (
	"context"
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/palantir/stacktrace"

	"go-poc/utils/log"
)

const (
	SALES_CHANNEL_STEP = 2
	INVENTORY_STEP     = 2
)

type Schema_migration struct {
	Version int  `orm:"version" json:"version"`
	Dirty   bool `orm:"dirty" json:"dirty"`
}

func initMigration(ctx context.Context, db *sql.DB, service, driverType string) {
	var err error
	migrationsTable := service + "_schema_migrations"
	m := &migrate.Migrate{}
	if driverType == "mysql" {
		driver, _ := mysql.WithInstance(db, &mysql.Config{
			MigrationsTable: migrationsTable,
		})

		m, err = migrate.NewWithDatabaseInstance(
			"file://./service/"+service+"/migration",
			driverType,
			driver,
		)

		if err != nil {
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "init instance db migration error"))
		}
	} else if driverType == "postgres" {
		driver, _ := postgres.WithInstance(db, &postgres.Config{
			MigrationsTable: migrationsTable,
		})

		m, err = migrate.NewWithDatabaseInstance(
			"file://./service/"+service+"/migration",
			driverType,
			driver,
		)

		if err != nil {
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "init instance db migration error"))
		}
	} else {
		return
	}

	sqlStatement := "SELECT version, dirty FROM " + migrationsTable
	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "select schema migration error"))
		switch service {
		case "saleschannel":
			err = m.Steps(SALES_CHANNEL_STEP)
		case "inventory":
			err = m.Steps(INVENTORY_STEP)
		default:
			err = errors.New("service not found")
		}

		if err != nil {
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "init step migration error"))
		}
	} else {
		row := stmt.QueryRow()

		var schema Schema_migration
		err = row.Scan(&schema.Version, &schema.Dirty)
		deltaSchema := 0
		switch service {
		case "saleschannel":
			deltaSchema = SALES_CHANNEL_STEP - schema.Version
		case "inventory":
			deltaSchema = INVENTORY_STEP - schema.Version
		}

		switch err {
		case sql.ErrNoRows:
			switch service {
			case "saleschannel":
				err = m.Steps(SALES_CHANNEL_STEP)
			case "inventory":
				err = m.Steps(INVENTORY_STEP)
			default:
				err = errors.New("service not found")
			}

			if err != nil {
				log.WithContext(ctx).Error(stacktrace.Propagate(err, "init step migration error"))
			}
		case nil:
			err = m.Steps(deltaSchema)
			if err != nil {
				if err.Error() != "no change" {
					log.WithContext(ctx).Error(stacktrace.Propagate(err, "init step migration error"))
				}
			}
		default:
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "migration error"))
		}
	}
}
