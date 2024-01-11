package storage

import (
	"gorm.io/gorm"
	"tinkoff-investment-bot/internal/model/database"
)

type UserPostgresStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) *UserPostgresStorage {
	return &UserPostgresStorage{db: db}
}

func (u *UserPostgresStorage) AddUser(user *database.User) error {
	result := u.db.FirstOrCreate(user, database.User{TelegramID: user.TelegramID})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *UserPostgresStorage) GetUserByTelegramChatID(telegramChatID int64) database.User {
	var user database.User
	u.db.Table("users").Select("*").Where("telegram_id = ?", telegramChatID).Scan(&user)
	return user
}
