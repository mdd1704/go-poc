package adapter

import (
	"database/sql"

	"github.com/pkg/errors"

	"go-poc/service/inventory/repository/adapter/location"
	"go-poc/service/inventory/repository/adapter/sourcing"
	"go-poc/service/inventory/repository/port"
	"go-poc/utils"
)

type postgresRegistry struct {
	db         *sql.DB
	dbexecutor utils.DBExecutor
}

func NewPostgres(db *sql.DB) port.MainRepository {
	return postgresRegistry{
		db: db,
	}
}

func (r postgresRegistry) Location() port.LocationMainRepository {
	if r.dbexecutor != nil {
		return location.NewPostgresRepository(r.dbexecutor)
	}
	return location.NewPostgresRepository(r.db)
}

func (r postgresRegistry) Sourcing() port.SourcingMainRepository {
	if r.dbexecutor != nil {
		return sourcing.NewMySQLRepository(r.dbexecutor)
	}
	return sourcing.NewMySQLRepository(r.db)
}

func (r postgresRegistry) DoInTransaction(txFunc port.InTransaction) (out interface{}, err error) {
	var tx *sql.Tx
	registry := r
	if r.dbexecutor == nil {
		tx, err = r.db.Begin()
		if err != nil {
			return
		}
		defer func() {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				switch x := p.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					// Fallback err (per specs, error strings should be lowercase w/o punctuation
					err = errors.New("unknown panic")
				}
			} else if err != nil {
				xerr := tx.Rollback() // err is non-nil; don't change it
				if xerr != nil {
					err = errors.Wrap(err, xerr.Error())
				}
			} else {
				err = tx.Commit() // err is nil; if Commit returns error update err
			}
		}()
		registry = postgresRegistry{
			db:         r.db,
			dbexecutor: tx,
		}
	}
	out, err = txFunc(registry)
	if err != nil {
		if out != nil {
			return out, err
		}

		return nil, err
	}
	return
}
