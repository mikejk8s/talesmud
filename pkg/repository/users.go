package repository

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/mikejk8s/talesmud/pkg/db"
	e "github.com/mikejk8s/talesmud/pkg/entities"
)

//UsersRepository ...
type UsersRepository interface {
	FindByRefID(id string) (*e.User, error)
	FindByID(id string) (*e.User, error)
	FindAll() ([]*e.User, error)
	FindAllOnline() ([]*e.User, error)
	Create(user *e.User) (*e.User, error)
	Import(user *e.User) (*e.User, error)
	Update(id string, user *e.User) error
	Delete(id string) error
	Drop() error
}

//--- Implementations

type usersRepo struct {
	*GenericRepo
}

//NewMongoDBUsersRepository creates a new mongodb partiesrep
func NewMongoDBUsersRepository(db *db.Client) UsersRepository {
	ur := &usersRepo{
		GenericRepo: &GenericRepo{
			db:         db,
			collection: "users",
			generator: func() interface{} {
				return &e.User{}
			},
		},
	}

	ur.CreateIndex()
	return ur
}

// Drop ...
func (pr *usersRepo) Drop() error {
	return pr.GenericRepo.DropCollection()
}

// Import ...
func (pr *usersRepo) Import(user *e.User) (*e.User, error) {
	result, err := pr.GenericRepo.Store(user)
	return result.(*e.User), err
}

func (pr *usersRepo) Create(user *e.User) (*e.User, error) {
	user.Entity = e.NewEntity()
	return pr.Import(user)
}

func (pr *usersRepo) FindAll() ([]*e.User, error) {
	results := make([]*e.User, 0)
	if err := pr.GenericRepo.FindAll(func(elem interface{}) {
		results = append(results, elem.(*e.User))
	}); err != nil {
		return nil, err
	}
	return results, nil
}

func (pr *usersRepo) FindAllOnline() ([]*e.User, error) {
	results := make([]*e.User, 0)
	if err := pr.GenericRepo.FindAllWithParam(
		db.NewQueryParams(db.QueryParam{Key: "isOnline", Value: true}),
		func(elem interface{}) {
			results = append(results, elem.(*e.User))
		}); err != nil {
		return nil, err
	}
	return results, nil
}

func (pr *usersRepo) Update(refID string, user *e.User) error {
	return pr.GenericRepo.UpdateByField(user, "refid", refID)
}
func (pr *usersRepo) FindByID(id string) (*e.User, error) {

	if id == "" {
		log.Error("Users::FindByID - ID is empty")
		return nil, errors.New("Empty ID")
	}

	result, err := pr.GenericRepo.FindByID(id)

	if user, ok := result.(*e.User); ok {
		return user, nil
	}
	return nil, err //ors.New("result is not a User")
}
func (pr *usersRepo) FindByRefID(refID string) (*e.User, error) {

	if refID == "" {
		log.Error("Users::FindByRefID - refID is empty")
		return nil, errors.New("Empty refID")
	}

	result, err := pr.GenericRepo.FindByField("refid", refID)

	if user, ok := result.(*e.User); ok {
		return user, nil
	}
	return nil, err //ors.New("result is not a User")
}

func (pr *usersRepo) Delete(id string) error {
	return pr.GenericRepo.Delete(id)
}
