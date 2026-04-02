package util

import (
	"errors"
	"unicode"
)

// ValidatePassword to check security same as in frontend
// Minimum of 8 characters
// Minmum of 1 lowercase character
// Minmum of 1 uppercase character
// Minimum of 1 Symbil
func ValidatePassword(password string) error {
	var hasUpper, hasLower, hasSymbol bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			// IsPunct to check reading marks like ( ! ? , .)
			// IsSymbol to check mathematic symbol like ( $ + =)
			hasSymbol = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasSymbol {
		return errors.New("password must contain at least one symbol/special character")
	}

	return nil
}
