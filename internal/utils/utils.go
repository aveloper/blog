package utils

import (
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

func ValidatePassword(pwd string) bool {
	tests := []string{"[a-z]", "[A-Z]", "[0-9]", "[^\\d\\w]"}
	for _, test := range tests {
		t, _ := regexp.MatchString(test, pwd)
		if !t {
			return false
		}
	}

	return true
}

func HashAndSalt(pwd []byte) (string, error) {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}    // GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

func ComparePasswords(hashedPwd string, plainPwd []byte) error {
	// Since we'll be getting the hashed password from the DB it
	// will be a string, so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return err
	}

	return nil
}

func ValidateAndHashPassword(pwd string) (string, bool, error ) {
	ok := ValidatePassword(pwd)
	if !ok {
		return "", false, nil
	}

	hash, err := HashAndSalt([]byte(pwd))
	if err != nil {
		return "", false ,err
	}

	return  hash, true ,nil
}