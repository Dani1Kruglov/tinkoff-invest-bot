package storage

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"sync"
	printbot "tinkoff-investment-bot/internal/bot/print"
	"tinkoff-investment-bot/internal/model/database"
	"tinkoff-investment-bot/internal/model/tracker"
	"tinkoff-investment-bot/internal/services/marketdata"
)

type SharePostgresStorage struct {
	db *gorm.DB
}

type Result struct {
	UID        string  `gorm:"type:varchar(200)"`
	Ticker     string  `gorm:"column:ticker;type:varchar(10)"`
	Name       string  `gorm:"column:name;type:varchar(100)"`
	TelegramID string  `gorm:"column:telegram_id;type:varchar(100)"`
	LastPrice  float32 `gorm:"column:last_price;type:float"`
}

func NewShareStorage(db *gorm.DB) *SharePostgresStorage {
	return &SharePostgresStorage{db: db}
}

func (s *SharePostgresStorage) AddShare(share *database.Share, userID uint, price float32) error {
	result := s.db.FirstOrCreate(share, database.Share{UID: share.UID})
	if result.Error != nil {
		return result.Error
	}

	result = s.db.FirstOrCreate(&database.UserFavoriteShare{
		UserID:    userID,
		ShareID:   share.ID,
		LastPrice: price,
	}, database.UserFavoriteShare{
		UserID:  userID,
		ShareID: share.ID,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *SharePostgresStorage) GetShares(tracker *tracker.Tracker) {
	var results []*Result
	s.db.Model(&database.Share{}).
		Select("shares.uid, shares.ticker, shares.Name, users.telegram_id, user_favorite_shares.last_price ").
		Joins("left join user_favorite_shares on user_favorite_shares.share_id = shares.id").
		Joins("left join users on users.id = user_favorite_shares.user_id").Scan(&results)
	var wg sync.WaitGroup
	for _, result := range results {
		wg.Add(1)
		share := result
		go func() {
			defer wg.Done()
			currentPrice, pastPrice, message, err := getSharesWithBigPriceChange(tracker, share)
			if err != nil {
				_ = fmt.Errorf(err.Error()) //
				return
			}
			if message != "" {
				printbot.PriceChange(currentPrice, pastPrice, message)
			}
		}()
	}
	wg.Wait()
}

func getSharesWithBigPriceChange(tracker *tracker.Tracker, share *Result) (float32, float32, string, error) {
	instrument, err := tracker.InstrumentsService.ShareByUid(share.UID)
	if err != nil {
		return 0, 0, "", err
	}
	currentPrice, err := marketdata.GetLastPriceByFigi(tracker, instrument.GetInstrument())
	if err != nil {
		return 0, 0, "", err
	}
	priceDifference := currentPrice * 100.00 / share.LastPrice
	fmt.Println(priceDifference)
	if priceDifference > 100.00 && priceDifference-100.00 > 5 {
		return currentPrice, share.LastPrice, "Цена выросла на " + strconv.FormatFloat(float64(priceDifference-100.00), 'f', -1, 32) + "%", nil
	} else if priceDifference < 100.00 && 100.00-priceDifference > 5 {
		return currentPrice, share.LastPrice, "Цена упала на " + strconv.FormatFloat(float64(100.00-priceDifference), 'f', -1, 32) + "%", nil
	}
	return 0, 0, "", nil
}
