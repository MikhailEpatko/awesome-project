package employee

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"idm/inner/common"
)

type Service struct {
	repo      Repo
	validator Validator
}

type Repo interface {
	FindById(id int64) (Entity, error)
	Save(entity Entity) (int64, error)
	BeginTransaction() (*sqlx.Tx, error)
	FindByNameTx(tx *sqlx.Tx, name string) (bool, error)
	SaveTx(tx *sqlx.Tx, employee Entity) (int64, error)
}

type Validator interface {
	Validate(request any) error
}

func NewService(
	repo Repo,
	validator Validator,
) *Service {
	return &Service{
		repo:      repo,
		validator: validator,
	}
}

func (svc *Service) FindById(id int64) (Response, error) {
	var entity, err = svc.repo.FindById(id)
	if err != nil {
		return Response{}, fmt.Errorf("error finding employee with id %d: %w", id, err)
	}
	return entity.toResponse(), nil
}

func (svc *Service) CreateEmployee(request CreateRequest) (int64, error) {
	var err = svc.validator.Validate(request)
	if err != nil {
		return 0, common.RequestValidationError{Message: err.Error()}
	}
	tx, err := svc.repo.BeginTransaction()
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	if err != nil {
		return 0, fmt.Errorf("error create employee: error creating transaction: %w", err)
	}
	isExist, err := svc.repo.FindByNameTx(tx, request.Name)
	if err != nil {
		return 0, fmt.Errorf("error finding employee by name: %s, %w", request.Name, err)
	}
	if isExist {
		return 0, common.AlreadyExistsError{Message: fmt.Sprintf("employee with name %s already exists", request.Name)}
	}
	newEmployeeId, err := svc.repo.SaveTx(tx, request.ToEntity())
	if err != nil {
		err = fmt.Errorf("error creating employee with name: %s %v", request.Name, err)
	}
	return newEmployeeId, err
}
