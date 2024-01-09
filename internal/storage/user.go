package storage

import (
	"gorm.io/gorm"
	"tinkoff-investment-bot/internal/model"
)

type UserPostgresStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) *UserPostgresStorage {
	return &UserPostgresStorage{db: db}
}

func (u *UserPostgresStorage) AddUser(user *model.User) error {
	result := u.db.FirstOrCreate(user, model.User{TelegramID: user.TelegramID, Token: user.Token})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *UserPostgresStorage) GetUserByTelegramID(telegramID string) model.User {
	var user model.User
	u.db.Table("users").Select("*").Where("telegram_id = ?", telegramID).Scan(&user)
	return user
}
