package repository

import (
	"auth-service/pkg/domain"
	"auth-service/pkg/repository/interfaces"

	"gorm.io/gorm"
)

type authDatabase struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) interfaces.AuthRepo {
	return &authDatabase{
		db: db,
	}
}

func (a *authDatabase) IsUserExist(email string) (bool, error) {

	res := a.db.Where("email = ?", email).Find(&domain.User{})
	if res.Error != nil {
		return false, res.Error
	}

	return res.RowsAffected != 0, nil
}

func (a *authDatabase) FindUserByEmail(email string) (user domain.User, err error) {

	err = a.db.Where("email = ?", email).Find(&user).Error

	return
}

func (a *authDatabase) SaveUser(user domain.User) (domain.User, error) {

	err := a.db.Create(&user).Error

	return user, err
}
