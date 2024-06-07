package main

import (
	"fmt"
	"strings"
	"unicode"
)

// checkPasswordStrength calculates the strength of the password and returns a score from 0 to 100
func checkPasswordStrength(password string) int {
	var score int

	// Criteria
	length := len(password)
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	specialChars := "!@#$%^&*()-_=+[]{}|;:,.<>?/`~"

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	// Length score
	if length >= 8 {
		score += 25
	} else {
		score += length * 3 // up to 21
	}

	// Uppercase letter score
	if hasUpper {
		score += 15
	}

	// Lowercase letter score
	if hasLower {
		score += 15
	}

	// Digit score
	if hasDigit {
		score += 20
	}

	// Special character score
	if hasSpecial {
		score += 25
	}

	// Ensure the score is capped at 100
	if score > 100 {
		score = 100
	}

	return score
}

func main() {
	passwords := []string{
		"password",
		"Password123",
		"P@ssw0rd!",
		"P@ssw0rd12345!",
		"P@ssw0rd12345!VeryLongPassword",
	}

	for _, password := range passwords {
		score := checkPasswordStrength(password)
		fmt.Printf("Password: %s, Strength Score: %d/100\n", password, score)
	}
}
