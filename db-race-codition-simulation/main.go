package main

import (
	"db-race-condition-simulation/config"
	"db-race-condition-simulation/module"
	"log"

	"gorm.io/gorm"
)

func main() {
	db, _ := config.InitDB()
	autoMigrate(db)
}

func autoMigrate(db *gorm.DB) {
	log.Println("Start migrate database")
	for _, table := range []interface{}{
		&module.User{},
		&module.Wallet{},
		&module.Transaction{},
	} {
		log.Printf("drop table %T", table)
		db.Migrator().DropTable(table)
		log.Printf("migrate table %T", table)
		err := db.AutoMigrate(table)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
