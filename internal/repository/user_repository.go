package repository

import (
	"github.com/azahir21/go-backend-boilerplate/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
    Create(user *model.User) error
    FindByUsername(username string) (*model.User, error)
    FindByEmail(email string) (*model.User, error)
    FindByID(id uint) (*model.User, error)
    Update(user *model.User) error
    Delete(id uint) error
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
    var user model.User
    err := r.db.Where("username = ?", username).First(&user).Error
    return &user, err
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
    var user model.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
    var user model.User
    err := r.db.First(&user, id).Error
    return &user, err
}

func (r *userRepository) Update(user *model.User) error {
    return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
    return r.db.Delete(&model.User{}, id).Error
}