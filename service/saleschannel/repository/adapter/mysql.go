package adapter

import (
	"database/sql"

	"github.com/pkg/errors"

	"go-poc/service/saleschannel/repository/adapter/channel"
	"go-poc/service/saleschannel/repository/port"
	"go-poc/utils"
)

type mysqlRegistry struct {
	db         *sql.DB
	dbexecutor utils.DBExecutor
}

func NewMySQL(db *sql.DB) port.MainRepository {
	return mysqlRegistry{
		db: db,
	}
}

func (r mysqlRegistry) Channel() port.ChannelMainRepository {
	if r.dbexecutor != nil {
		return channel.NewMySQLRepository(r.dbexecutor)
	}
	return channel.NewMySQLRepository(r.db)
}

func (r mysqlRegistry) DoInTransaction(txFunc port.InTransaction) (out interface{}, err error) {
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
		registry = mysqlRegistry{
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
