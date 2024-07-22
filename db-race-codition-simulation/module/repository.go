package module

import (
	"context"

	redis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	Register(name string) (*User, error)
	CreateWallet(user *User, initAmount float64) (*Wallet, error)
	GetWalletByID(id int) (*Wallet, error)
	TopupBalance(currentUser *User, wallet *Wallet, amount float64) error
	Transfer(fromUser, toUser *User, fromWallet, toWallet *Wallet, amount float64) error
	ValidateIdempotency(ctx context.Context, key string) (bool, error)
}

func NewRepository(db *gorm.DB, cache *redis.Client) Repository {
	return &repo{db, cache}
}

type repo struct {
	db    *gorm.DB
	cache *redis.Client
}

func (r *repo) Register(name string) (*User, error) {
	user := User{
		Name: name,
	}
	err := r.db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repo) CreateWallet(user *User, initAmount float64) (*Wallet, error) {
	var result *Wallet
	resultErr := r.db.Transaction(func(tx *gorm.DB) error {
		// create wallet
		wallet := Wallet{
			UserID:  user.ID,
			Balance: initAmount,
		}
		err := tx.Create(&wallet).Error
		if err != nil {
			return err
		}

		// create transaction
		err = tx.Create(&Transaction{
			UserID:      user.ID,
			WalletID:    wallet.ID,
			Flow:        TrxFlowCredit,
			Amount:      initAmount,
			Description: "init wallet",
		}).Error
		if err != nil {
			return err
		}

		// set result
		result = &wallet
		return nil
	})
	return result, resultErr
}

func (r *repo) GetWalletByID(id int) (*Wallet, error) {
	var result Wallet
	err := r.db.First(&result, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *repo) TopupBalance(currentUser *User, wallet *Wallet, amount float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// create transaction
		err := tx.Create(&Transaction{
			UserID:      currentUser.ID,
			WalletID:    wallet.ID,
			Flow:        TrxFlowCredit,
			Amount:      amount,
			Description: "topup",
		}).Error
		if err != nil {
			return err
		}

		// update balance
		return tx.Model(wallet).Update("balance", gorm.Expr("balance + ?", amount)).Error
	})
}

func (r *repo) Transfer(fromUser, toUser *User, fromWallet, toWallet *Wallet, amount float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// create debit transaction for fromUser
		err := tx.Create(&Transaction{
			UserID:      toUser.ID,
			WalletID:    fromWallet.ID,
			Flow:        TrxFlowDebit,
			Amount:      amount,
			Description: "transfer",
		}).Error
		if err != nil {
			return err
		}

		// create credit transaction for toUser
		err = tx.Create(&Transaction{
			UserID:      fromUser.ID,
			WalletID:    toWallet.ID,
			Flow:        TrxFlowCredit,
			Amount:      amount,
			Description: "transfer",
		}).Error
		if err != nil {
			return err
		}

		// update toUser balance
		err = tx.Model(toWallet).Update("balance", gorm.Expr("balance + ?", amount)).Error
		if err != nil {
			return err
		}

		// update fromUser balance
		return tx.Model(fromWallet).Update("balance", gorm.Expr("balance - ?", amount)).Error
	})
}

func (r *repo) ValidateIdempotency(ctx context.Context, key string) (bool, error) {
	tr := r.cache.HSetNX(ctx, "idempotency:"+key, "status", "started")
	if tr.Err() != nil {
		return false, tr.Err()
	}
	return !tr.Val(), nil
}
