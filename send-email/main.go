package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

const (
	smptHost     string = "SMTP_HOST"
	smptPort     string = "SMTP_PORT"
	smptUser     string = "SMTP_USER"
	smptPassword string = "SMTP_PASSWORD"
)

func main() {
	// read .env file
	env, err := godotenv.Read()
	if err != nil {
		log.Fatalln(err)
	}

	// set reader
	r := bufio.NewReader(os.Stdin)

	// create new message
	m := gomail.NewMessage()

	// set email header
	m.SetHeader("From", env[smptUser])
	m.SetHeader("To", input(r, "Enter email address: "))
	m.SetHeader("Subject", input(r, "Enter subject: "))

	// set email content
	m.SetBody("text/plain", multilineInput(r, "Enter email content"))

	// include attachment
	m.Attach(input(r, "Include file path to set as attachement: "))

	// set dialer
	port, err := strconv.Atoi(env[smptPort])
	if err != nil {
		log.Fatalln(err)
	}
	d := gomail.NewDialer(env[smptHost], port, env[smptUser], env[smptPassword])

	// Mengirim email
	if err := d.DialAndSend(m); err != nil {
		log.Fatalln(err)
	}

	log.Println("Email sent successfully with attachment!")
}

func input(reader *bufio.Reader, caption string) string {
	fmt.Print(caption)
	result, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}
	result = strings.TrimSpace(result)
	if result == "" {
		log.Fatalln("error, empty input")
	}
	return result
}

func multilineInput(reader *bufio.Reader, caption string) string {
	fmt.Println(caption, "(press Ctrl+D to finish):")

	var result string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		result += line
	}
	result = strings.TrimSpace(result)
	if result == "" {
		log.Fatalln("error, empty input")
	}
	return result
}
