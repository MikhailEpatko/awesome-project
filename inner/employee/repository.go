package employee

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(database *sqlx.DB) *Repository {
	return &Repository{db: database}
}

func (r *Repository) FindById(id int64) (employee Entity, err error) {
	err = r.db.Get(&employee, "select * from employee where id = $1", id)
	return employee, err
}

func (r *Repository) Save(employee Entity) (employeeId int64, err error) {
	err = r.db.Get(
		&employeeId,
		`insert into employee (name) values ($1) returning id`,
		employee.Name,
	)
	return employeeId, err
}

// transactional methods

func (r *Repository) BeginTransaction() (tx *sqlx.Tx, err error) {
	return r.db.Beginx()
}

func (r *Repository) FindByNameTx(tx *sqlx.Tx, name string) (isExists bool, err error) {
	err = tx.Get(
		&isExists,
		"select exists(select 1 from employee where name = $1)",
		name,
	)
	return isExists, err
}

func (r *Repository) SaveTx(tx *sqlx.Tx, employee Entity) (employeeId int64, err error) {
	err = tx.Get(
		&employeeId,
		`insert into employee (name) values ($1) returning id`,
		employee.Name,
	)
	return employeeId, err
}
