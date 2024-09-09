package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"runtime"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	TIME        = 1
	MEMORY      = 64 * 1024
	KEY_LENGTH  = 32
	SALT_LENGTH = 16
)

func HashPassword(password string) (string, error) {
	// Generate salt
	salt := make([]byte, SALT_LENGTH)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, TIME, MEMORY, uint8(runtime.NumCPU()), KEY_LENGTH)

	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	hashB64 := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("%s$%s", saltB64, hashB64), nil
}

func VerifyPassword(password string, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}

	hashed, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	hash := argon2.IDKey([]byte(password), salt, TIME, MEMORY, uint8(runtime.NumCPU()), KEY_LENGTH)

	if len(hash) != len(hashed) {
		return false, nil
	}

	for i := 0; i < len(hash); i++ {
		if hash[i] != hashed[i] {
			return false, nil
		}
	}

	return true, nil
}
