package auth

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ContextUser struct {
	UserUUID  uuid.UUID
	SpotifyID string
}

var AnonymousUser = &ContextUser{}

// SetHash calculates the hash of the given password and stores the hash
func SetHash(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

// Matches checks if the given plaintext Password matches the hash stored in the struct
func Matches(plaintextPassword string, hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
