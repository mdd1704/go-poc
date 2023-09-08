package location

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"

	"go-poc/service/inventory/model"
	"go-poc/service/inventory/repository/port"
	"go-poc/utils"
)

type postgresRepository struct {
	db utils.DBExecutor
}

func NewPostgresRepository(db utils.DBExecutor) port.LocationMainRepository {
	return &postgresRepository{
		db: db,
	}
}

func (repo *postgresRepository) Create(data *model.Location) error {
	fmt.Println("data: ", data)
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Insert("locations").Rows(
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
	fmt.Println("query: ", query)

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

func (repo *postgresRepository) Update(data *model.Location) error {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Update("locations").Set(
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

func (repo *postgresRepository) FindByID(id uuid.UUID) (result *model.Location, err error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From("locations")
	dataset = dataset.Where(goqu.Ex{"id": id})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, stacktrace.Propagate(err, "dataset error")
	}

	row := repo.db.QueryRow(query)
	result = &model.Location{}
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

func (repo *postgresRepository) FindByFilter(filter model.LocationFilter, lock bool) (result []*model.Location, err error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From("locations")
	dataset = repo.addFilter(dataset, filter)

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, stacktrace.Propagate(err, "dataset error")
	}
	fmt.Println("query: ", query)

	if lock {
		query += " FOR UPDATE"
	}

	res, err := repo.db.Query(query)
	if err != nil {
		return nil, stacktrace.Propagate(err, "query error")
	}

	locations := []*model.Location{}
	for res.Next() {
		item := &model.Location{}
		err := res.Scan(
			&item.ID,
			&item.Code,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		locations = append(locations, item)
	}

	return locations, nil
}

func (repo *postgresRepository) FindPage(filter model.LocationFilter, offset, limit int64) (result []*model.Location, err error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From("locations")
	dataset = repo.addFilter(dataset, filter).Offset(uint(offset)).Limit(uint(limit))

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, stacktrace.Propagate(err, "dataset error")
	}

	res, err := repo.db.Query(query)
	if err != nil {
		return nil, stacktrace.Propagate(err, "query error")
	}

	locations := []*model.Location{}
	for res.Next() {
		item := &model.Location{}
		err := res.Scan(
			&item.ID,
			&item.Code,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		locations = append(locations, item)
	}

	return locations, nil
}

func (repo *postgresRepository) FindTotalByFilter(filter model.LocationFilter) (total int64, err error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From("locations")
	dataset = dataset.Select(goqu.COUNT("*"))
	dataset = repo.addFilter(dataset, filter)

	query, _, err := dataset.ToSQL()
	if err != nil {
		return 0, stacktrace.Propagate(err, "query error")
	}

	err = repo.db.QueryRow(query).Scan(&total)
	if err != nil {
		return 0, stacktrace.Propagate(err, "scan error")
	}

	return total, nil
}

func (repo *postgresRepository) Delete(filter model.LocationFilter) error {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Delete("locations")
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

func (repo *postgresRepository) addFilter(dataset *goqu.SelectDataset, filter model.LocationFilter) *goqu.SelectDataset {
	if len(filter.IDs) != 0 {
		dataset = dataset.Where(goqu.Ex{"id": filter.IDs})
	}

	if len(filter.Codes) != 0 {
		dataset = dataset.Where(goqu.Ex{"code": filter.Codes})
	}

	return dataset
}
