package repository

import (
	"errors"

	"github.com/gendutski/my-playground/dockerexample/core/entity"
	"gorm.io/gorm"
)

const defaultLimit int = 10

type Repository interface {
	Get(offset, limit int) ([]*entity.User, error)
	Set(payload *entity.User) error
}

type repo struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Repository {
	return &repo{db}
}

func (r *repo) Get(offset, limit int) ([]*entity.User, error) {
	if limit == 0 {
		limit = defaultLimit
	}
	var result []*entity.User
	err := r.db.Offset(offset).Limit(limit).Find(&result).Error
	return result, err
}

func (r *repo) Set(payload *entity.User) error {
	if payload.Gender == entity.InvalidGender {
		return errors.New("invalid gender")
	}
	return r.db.Save(payload).Error
}
