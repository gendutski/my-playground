package main

import (
	"bufio"
	"log"
	"os"

	"github.com/gendutski/my-playground/mysql-migration/utils/migrator"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	result, err := migrator.Run(reader)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("success")
	log.Println(result)
}
