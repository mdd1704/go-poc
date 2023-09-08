package channel

import (
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"

	"go-poc/service/saleschannel/model"
	"go-poc/service/saleschannel/repository/port"
	"go-poc/utils"
)

type mysqlRepository struct {
	db utils.DBExecutor
}

func NewMySQLRepository(db utils.DBExecutor) port.ChannelMainRepository {
	return &mysqlRepository{
		db: db,
	}
}

func (repo *mysqlRepository) Create(data *model.Channel) error {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.Insert("channels").Rows(
		goqu.Record{
			"id":         data.ID,
			"code":       data.Code,
			"created_at": data.CreatedAt,
			"updated_at": data.UpdatedAt,
		},
	)

	query, _, err := dataset.ToSQL()
	if err != nil {
		return stacktrace.Propagate(err, "dataset error")
	}

	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return stacktrace.Propagate(err, "prepare error")
	}

	_, err = stmt.Exec()
	if err != nil {
		return stacktrace.Propagate(err, "exec error")
	}

	return nil
}

func (repo *mysqlRepository) Update(data *model.Channel) error {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.Update("channels").Set(
		goqu.Record{
			"code":       data.Code,
			"updated_at": data.UpdatedAt,
		},
	)
	dataset = dataset.Where(goqu.Ex{"id": data.ID})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return stacktrace.Propagate(err, "dataset error")
	}

	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return stacktrace.Propagate(err, "prepare error")
	}

	_, err = stmt.Exec()
	if err != nil {
		return stacktrace.Propagate(err, "exec error")
	}

	return nil
}

func (repo *mysqlRepository) FindByID(id uuid.UUID) (result *model.Channel, err error) {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("channels")
	dataset = dataset.Where(goqu.Ex{"id": id})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, stacktrace.Propagate(err, "dataset error")
	}

	row := repo.db.QueryRow(query)
	result = &model.Channel{}
	err = row.Scan(
		&result.ID,
		&result.Code,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, stacktrace.Propagate(err, "scan error")
	}

	return result, nil
}

func (repo *mysqlRepository) FindByFilter(filter model.ChannelFilter, lock bool) (result []*model.Channel, err error) {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("channels")
	dataset = repo.addFilter(dataset, filter)

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, stacktrace.Propagate(err, "dataset error")
	}

	if lock {
		query += " FOR UPDATE"
	}

	res, err := repo.db.Query(query)
	if err != nil {
		return nil, stacktrace.Propagate(err, "query error")
	}

	channels := []*model.Channel{}
	for res.Next() {
		item := &model.Channel{}
		err := res.Scan(
			&item.ID,
			&item.Code,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		channels = append(channels, item)
	}

	return channels, nil
}

func (repo *mysqlRepository) FindPage(filter model.ChannelFilter, offset, limit int64) (result []*model.Channel, err error) {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("channels")
	dataset = repo.addFilter(dataset, filter).Offset(uint(offset)).Limit(uint(limit))

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, stacktrace.Propagate(err, "dataset error")
	}

	res, err := repo.db.Query(query)
	if err != nil {
		return nil, stacktrace.Propagate(err, "query error")
	}

	channels := []*model.Channel{}
	for res.Next() {
		item := &model.Channel{}
		err := res.Scan(
			&item.ID,
			&item.Code,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		channels = append(channels, item)
	}

	return channels, nil
}

func (repo *mysqlRepository) FindTotalByFilter(filter model.ChannelFilter) (total int64, err error) {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("channels")
	dataset = dataset.Select(goqu.COUNT("*"))
	dataset = repo.addFilter(dataset, filter)

	query, _, err := dataset.ToSQL()
	if err != nil {
		return 0, stacktrace.Propagate(err, "dataset error")
	}

	err = repo.db.QueryRow(query).Scan(&total)
	if err != nil {
		return 0, stacktrace.Propagate(err, "scan error")
	}

	return total, nil
}

func (repo *mysqlRepository) Delete(filter model.ChannelFilter) error {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.Delete("channels")
	dataset = dataset.Where(goqu.Ex{"id": filter.IDs})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return stacktrace.Propagate(err, "dataset error")
	}

	_, err = repo.db.Query(query)
	if err != nil {
		return stacktrace.Propagate(err, "query error")
	}

	return nil
}

func (repo *mysqlRepository) addFilter(dataset *goqu.SelectDataset, filter model.ChannelFilter) *goqu.SelectDataset {
	if len(filter.IDs) != 0 {
		dataset = dataset.Where(goqu.Ex{"id": filter.IDs})
	}

	if len(filter.Codes) != 0 {
		dataset = dataset.Where(goqu.Ex{"code": filter.Codes})
	}

	return dataset
}
