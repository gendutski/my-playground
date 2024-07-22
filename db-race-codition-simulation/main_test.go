package main

import (
	"context"
	"db-race-condition-simulation/config"
	"db-race-condition-simulation/module"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type result struct {
	Username string
	Amount   float64
	Error    error
}

func Test_Transfer(t *testing.T) {
	db, cache := config.InitDB()
	autoMigrate(db, t)

	// init module
	repo := module.NewRepository(db, cache)
	service := module.NewService(repo)

	// register users
	err := service.RegisterUsers()
	if err != nil {
		t.Fatalf("error call service.RegisterUsers: %s", err.Error())
	}
	// create users wallet
	err = service.CreateWalletPerUsers()
	if err != nil {
		t.Fatalf("error call service.CreateWalletPerUsers: %s", err.Error())
	}

	var transferAmount float64 = 100

	// topup user1
	var topup float64 = transferAmount * 1000
	err = service.TopupBalance(module.User1, topup)
	if err != nil {
		t.Fatalf("error call service.TopupBalance: %s", err.Error())
	}

	// set channel & context
	resultCh := make(chan result, 1000)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		key := uuid.New().String()

		// transfer user 1 -> user 2
		wg.Add(1)
		go func(ctx context.Context, key string, ch chan<- result) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				t.Log("timeout")
			default:
				err := service.Transfer(ctx, module.User1, module.User2, key, transferAmount)
				ch <- result{
					Username: module.User1,
					Amount:   transferAmount,
					Error:    err,
				}
			}
		}(ctx, key, resultCh)

		// duplicate transfer user 1 -> user 2
		wg.Add(1)
		go func(ctx context.Context, key string, ch chan<- result) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				t.Log("timeout")
			default:
				err := service.Transfer(ctx, module.User1, module.User2, key, transferAmount)
				ch <- result{
					Username: module.User1,
					Amount:   transferAmount,
					Error:    err,
				}
			}
		}(ctx, key, resultCh)

		// transfer user 3 -> user 2
		wg.Add(1)
		go func(ctx context.Context, ch chan<- result) {
			key := uuid.New().String()
			defer wg.Done()
			select {
			case <-ctx.Done():
				t.Log("timeout")
			default:
				err := service.Transfer(ctx, module.User3, module.User2, key, transferAmount)
				ch <- result{
					Username: module.User3,
					Amount:   transferAmount,
					Error:    err,
				}
			}
		}(ctx, resultCh)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var total1, total3 float64
	for data := range resultCh {
		if data.Error == nil {
			if data.Username == module.User1 {
				total1 += data.Amount
			} else {
				total3 += data.Amount
			}
		} else {
			if data.Username == module.User1 {
				assert.Equal(t, "duplicate request", data.Error.Error())
			}
		}
	}
	assert.Equal(t, topup, total1)
	assert.Equal(t, float64(10000), total3)
}

func autoMigrate(db *gorm.DB, t *testing.T) {
	t.Log("Start migrate database")
	for _, table := range []interface{}{
		&module.User{},
		&module.Wallet{},
		&module.Transaction{},
	} {
		t.Logf("drop table %T", table)
		db.Migrator().DropTable(table)
		t.Logf("migrate table %T", table)
		err := db.AutoMigrate(table)
		if err != nil {
			t.Fatal(err)
		}
	}
}
