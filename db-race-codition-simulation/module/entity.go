package module

import (
	"time"
)

type TrxFlow string

const (
	TrxFlowDebit  TrxFlow = "debit"
	TrxFlowCredit TrxFlow = "credit"
	User1         string  = "User no 1"
	User2         string  = "User no 2"
	User3         string  = "User no 3"
)

type UserWallet struct {
	User   *User
	Wallet *Wallet
}

type User struct {
	ID        int       `gorm:"primaryKey"`
	Name      string    `gorm:"size:30;not null;unique"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null"`
}

type Wallet struct {
	ID      int     `gorm:"primaryKey"`
	UserID  int     `gorm:"index;not null"`
	Balance float64 `gorm:"not null;type:decimal(10,2);default:0;check:balance>=0"`
}

type Transaction struct {
	UserID      int       `gorm:"index:idx_transacation;not null"`
	WalletID    int       `gorm:"index:idx_transacation;not null"`
	Flow        TrxFlow   `gorm:"index:idx_transacation;not null;type:enum('debit', 'credit')"`
	Amount      float64   `gorm:"not null;type:decimal(10,2);default:0"`
	Description string    `gorm:"not null;size:100"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null"`
}
