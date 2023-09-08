package sourcing

import (
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"

	"go-poc/service/inventory/model"
	"go-poc/service/inventory/repository/port"
	"go-poc/utils"
)

type mysqlRepository struct {
	db utils.DBExecutor
}

func NewMySQLRepository(db utils.DBExecutor) port.SourcingMainRepository {
	return &mysqlRepository{
		db: db,
	}
}

func (repo *mysqlRepository) Create(data *model.Sourcing) error {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.Insert("sourcings").Rows(
		goqu.Record{
			"id":           data.ID,
			"sku":          data.SKU,
			"qty_total":    data.QtyTotal,
			"qty_reserved": data.QtyReserved,
			"qty_saleable": data.QtySaleable,
			"created_at":   data.CreatedAt,
			"updated_at":   data.UpdatedAt,
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

func (repo *mysqlRepository) Update(data *model.Sourcing) error {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.Update("sourcings").Set(
		goqu.Record{
			"qty_total":    data.QtyTotal,
			"qty_reserved": data.QtyReserved,
			"qty_saleable": data.QtySaleable,
			"updated_at":   data.UpdatedAt,
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

func (repo *mysqlRepository) FindByID(id uuid.UUID) (result *model.Sourcing, err error) {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("sourcings")
	dataset = dataset.Where(goqu.Ex{"id": id})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, stacktrace.Propagate(err, "dataset error")
	}

	row := repo.db.QueryRow(query)
	result = &model.Sourcing{}
	err = row.Scan(
		&result.ID,
		&result.SKU,
		&result.QtyTotal,
		&result.QtyReserved,
		&result.QtySaleable,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, stacktrace.Propagate(err, "scan error")
	}

	return result, nil
}

func (repo *mysqlRepository) FindByFilter(filter model.SourcingFilter, lock bool) (result []*model.Sourcing, err error) {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("sourcings")
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

	Sourcings := []*model.Sourcing{}
	for res.Next() {
		item := &model.Sourcing{}
		err := res.Scan(
			&item.ID,
			&item.SKU,
			&item.QtyTotal,
			&item.QtyReserved,
			&item.QtySaleable,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		Sourcings = append(Sourcings, item)
	}

	return Sourcings, nil
}

func (repo *mysqlRepository) FindPage(filter model.SourcingFilter, offset, limit int64) (result []*model.Sourcing, err error) {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("sourcings")
	dataset = repo.addFilter(dataset, filter).Offset(uint(offset)).Limit(uint(limit))

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, stacktrace.Propagate(err, "dataset error")
	}

	res, err := repo.db.Query(query)
	if err != nil {
		return nil, stacktrace.Propagate(err, "query error")
	}

	Sourcings := []*model.Sourcing{}
	for res.Next() {
		item := &model.Sourcing{}
		err := res.Scan(
			&item.ID,
			&item.SKU,
			&item.QtyTotal,
			&item.QtyReserved,
			&item.QtySaleable,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		Sourcings = append(Sourcings, item)
	}

	return Sourcings, nil
}

func (repo *mysqlRepository) FindTotalByFilter(filter model.SourcingFilter) (total int64, err error) {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("sourcings")
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

func (repo *mysqlRepository) Delete(filter model.SourcingFilter) error {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.Delete("sourcings")
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

func (repo *mysqlRepository) addFilter(dataset *goqu.SelectDataset, filter model.SourcingFilter) *goqu.SelectDataset {
	if len(filter.IDs) != 0 {
		dataset = dataset.Where(goqu.Ex{"id": filter.IDs})
	}

	if len(filter.SKUs) != 0 {
		dataset = dataset.Where(goqu.Ex{"sku": filter.SKUs})
	}

	return dataset
}
