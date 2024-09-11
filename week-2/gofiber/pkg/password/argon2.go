package password

import (
	"bytes"
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

func generateSalt() ([]byte, error) {
	salt := make([]byte, SALT_LENGTH)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

func Hash(password string) (string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", err
	}

	hashed := argon2.IDKey([]byte(password), salt, TIME, MEMORY, uint8(runtime.NumCPU()), KEY_LENGTH)

	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	hashB64 := base64.RawStdEncoding.EncodeToString(hashed)

	return fmt.Sprintf("%s$%s", saltB64, hashB64), nil
}

func Verify(password, hashedPassword string) (bool, error) {
	parts := strings.Split(hashedPassword, "$")

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}

	hashed, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	hash := argon2.IDKey([]byte(password), salt, TIME, MEMORY, uint8(runtime.NumCPU()), KEY_LENGTH)

	if bytes.Equal(hash, hashed) {
		return false, nil
	} else {
		return true, nil
	}
}
