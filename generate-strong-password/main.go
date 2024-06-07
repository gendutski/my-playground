package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	passwordLength = 16
	letters        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits         = "0123456789"
	specialChars   = "!@#$%^&*()-_=+[]{}|;:,.<>?/`~"
)

// generateRandomChar generates a random character from the given charset
func generateRandomChar(charset string) (byte, error) {
	max := big.NewInt(int64(len(charset)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return charset[n.Int64()], nil
}

// generateStrongPassword generates a strong password of the specified length
func generateStrongPassword(length int) (string, error) {
	if length < 8 {
		return "", fmt.Errorf("password length should be at least 8 characters")
	}

	var password strings.Builder

	// Ensure password contains at least one character from each set
	charsets := []string{letters, digits, specialChars}

	for _, charset := range charsets {
		char, err := generateRandomChar(charset)
		if err != nil {
			return "", err
		}
		password.WriteByte(char)
	}

	// Fill the rest of the password length with random characters from all sets
	allChars := letters + digits + specialChars
	for password.Len() < length {
		char, err := generateRandomChar(allChars)
		if err != nil {
			return "", err
		}
		password.WriteByte(char)
	}

	// Convert the password to a slice to shuffle the characters
	passwordSlice := []byte(password.String())
	for i := range passwordSlice {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(len(passwordSlice))))
		if err != nil {
			return "", err
		}
		passwordSlice[i], passwordSlice[j.Int64()] = passwordSlice[j.Int64()], passwordSlice[i]
	}

	return string(passwordSlice), nil
}

func main() {
	password, err := generateStrongPassword(passwordLength)
	if err != nil {
		fmt.Printf("Error generating password: %v\n", err)
		return
	}

	fmt.Printf("Generated password: %s\n", password)
}
