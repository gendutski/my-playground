package terminal

import (
	"bufio"
	"fmt"
	"strings"
)

func ReadString(prompt, defaultValue string, reader *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	result, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	result = trim(result, defaultValue)
	return result, nil
}

func Confirm(prompt string, reader *bufio.Reader) bool {
	fmt.Print(prompt)
	result, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("error confirm:", err.Error())
	}
	result = trim(result, "")
	if result == "y" || result == "Y" || strings.EqualFold(result, "yes") {
		return true
	}
	return false
}

func trim(value, defaultValue string) string {
	value = strings.TrimSpace(value)
	if value == "" && defaultValue != "" {
		return defaultValue
	}
	return value
}
