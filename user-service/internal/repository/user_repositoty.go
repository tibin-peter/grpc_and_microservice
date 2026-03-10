package repository

import (
	"grpc_and_microservice/user-service/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user *model.User) error
	FindUserByEmail(email string) (*model.User, error)
	FindUserByID(id string) (*model.User, error)
}

type PostgreSQL struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &PostgreSQL{db: db}
}

func (r *PostgreSQL) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *PostgreSQL) FindUserByEmail(email string) (*model.User, error) {

	var user model.User

	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *PostgreSQL) FindUserByID(id string) (*model.User, error) {

	var user model.User

	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}