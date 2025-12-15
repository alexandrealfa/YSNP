package entity

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type ID = uuid.UUID

func NewId() ID {
	return ID(uuid.New())
}

func NewPassword(passwordLength int, specialCharacters bool) string {
	lowerCase := "abcdefghijklmnopqrstuvwxyz" // lower
	upperCase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" // upper
	numbers := "0123456789"                   // numbers
	specialChar := "!@#$%^&*()_-+={}[/?]"     // special

	password := ""
	numbersToRand := 3

	if specialCharacters {
		numbersToRand += 1
	}

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	for n := 0; n < passwordLength; n++ {
		randNum := rng.Intn(numbersToRand)

		switch randNum {
		case 0:
			randCharNum := rng.Intn(len(lowerCase))
			password += string(lowerCase[randCharNum])
		case 1:
			randCharNum := rng.Intn(len(upperCase))
			password += string(upperCase[randCharNum])
		case 2:
			randCharNum := rng.Intn(len(numbers))
			password += string(numbers[randCharNum])
		case 3:
			randCharNum := rng.Intn(len(specialChar))
			password += string(specialChar[randCharNum])
		}
	}

	return password
}
