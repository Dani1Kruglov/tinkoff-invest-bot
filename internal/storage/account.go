package storage

import (
	"errors"
	"gorm.io/gorm"
	"tinkoff-investment-bot/internal/model"
)

type AccountPostgresStorage struct {
	db *gorm.DB
}

func NewAccountStorage(db *gorm.DB) *AccountPostgresStorage {
	return &AccountPostgresStorage{db: db}
}

func (a *AccountPostgresStorage) AddAccount(account *model.Account, userID uint) error {
	if userID != 0 {
		account.UserID = userID
		result := a.db.FirstOrCreate(account, model.Account{AccountID: account.AccountID, UserID: userID})
		if result.Error != nil {
			return result.Error
		}
		return nil
	} else {
		return errors.New("error select user by telegram id")
	}
}

func (a *AccountPostgresStorage) GetAccountIDByTelegramID(telegramID string) model.Account {
	var account model.Account
	a.db.Table("accounts").Select("*").Joins("left join users on accounts.user_id = users.id").Where("users.telegram_id = ?", telegramID).Scan(&account)
	return account
}
