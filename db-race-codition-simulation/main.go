package main

// import (
// 	"context"
// 	"db-race-condition-simulation/config"
// 	"db-race-condition-simulation/module"
// 	"log"
// 	"os"
// 	"sync"
// 	"time"

// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )

// type result struct {
// 	Username string
// 	Amount   float64
// 	Error    error
// }

// func main() {
// 	db, cache := config.InitDB()
// 	autoMigrate(db)

// 	// create log file
// 	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY, 0666)
// 	if err != nil {
// 		log.Fatalf("Failed to open log file: %v", err)
// 	}
// 	defer file.Close()

// 	// init module
// 	repo := module.NewRepository(db, cache)
// 	service := module.NewService(repo)

// 	// register users
// 	err = service.RegisterUsers()
// 	if err != nil {
// 		log.Fatalf("error call service.RegisterUsers: %s", err.Error())
// 	}
// 	// create users wallet
// 	err = service.CreateWalletPerUsers()
// 	if err != nil {
// 		log.Fatalf("error call service.CreateWalletPerUsers: %s", err.Error())
// 	}

// 	var transferAmount float64 = 100

// 	// topup user1
// 	err = service.TopupBalance(module.User1, transferAmount*1000)
// 	if err != nil {
// 		log.Fatalf("error call service.TopupBalance: %s", err.Error())
// 	}

// 	// set channel & context
// 	resultCh := make(chan result, 100)
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()

// 	var wg sync.WaitGroup
// 	for i := 0; i < 1000; i++ {
// 		key := uuid.New().String()

// 		// transfer user 1 -> user 2
// 		wg.Add(1)
// 		go func(ctx context.Context, key string, ch chan<- result) {
// 			defer wg.Done()
// 			err := service.Transfer(ctx, module.User1, module.User2, key, transferAmount)
// 			ch <- result{
// 				Username: module.User1,
// 				Amount:   transferAmount,
// 				Error:    err,
// 			}
// 		}(ctx, key, resultCh)

// 		// duplicate transfer user 1 -> user 2
// 		wg.Add(1)
// 		go func(ctx context.Context, key string, ch chan<- result) {
// 			defer wg.Done()
// 			err := service.Transfer(ctx, module.User1, module.User2, key, transferAmount)
// 			ch <- result{
// 				Username: module.User1,
// 				Amount:   transferAmount,
// 				Error:    err,
// 			}
// 		}(ctx, key, resultCh)

// 		// transfer user 3 -> user 2
// 		wg.Add(1)
// 		go func(ctx context.Context, ch chan<- result) {
// 			key := uuid.New().String()

// 			defer wg.Done()
// 			err := service.Transfer(ctx, module.User3, module.User2, key, transferAmount)
// 			ch <- result{
// 				Username: module.User3,
// 				Amount:   transferAmount,
// 				Error:    err,
// 			}
// 		}(ctx, resultCh)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(resultCh)
// 	}()

// 	var total1, total3 float64
// 	log.SetOutput(file)
// 	for data := range resultCh {
// 		if data.Error != nil {
// 			log.Println(err)
// 		} else if data.Username == module.User1 {
// 			total1 += data.Amount
// 		} else {
// 			total3 += data.Amount
// 		}
// 	}

// 	log.SetOutput(os.Stdin)
// 	log.Printf("total transfered from %s = %f", module.User1, total1)
// 	log.Printf("total transfered from %s = %f", module.User3, total3)
// }

// func autoMigrate(db *gorm.DB) {
// 	log.Println("Start migrate database")
// 	for _, table := range []interface{}{
// 		&module.User{},
// 		&module.Wallet{},
// 		&module.Transaction{},
// 	} {
// 		log.Printf("drop table %T", table)
// 		db.Migrator().DropTable(table)
// 		log.Printf("migrate table %T", table)
// 		err := db.AutoMigrate(table)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 	}
// }
