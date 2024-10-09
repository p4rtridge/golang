package helpers

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Hasher interface {
	HashPassword(password string) (string, error)
	CompareHash(hashedPassword, password string) (bool, error)
}

type hasher struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func NewHasher(memory, iterations, saltLength, keyLength uint32, parallelism uint8) Hasher {
	return &hasher{
		memory,
		iterations,
		parallelism,
		saltLength,
		keyLength,
	}
}

func (r *hasher) RandomBytes() ([]byte, error) {
	salt := make([]byte, r.saltLength)

	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

func (r *hasher) EncodeHash(password, salt []byte) string {
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(password)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, r.memory, r.iterations, r.parallelism, b64Salt, b64Hash)
}

func (r *hasher) DecodeHash(encodedHash string) (*hasher, []byte, []byte, error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	var p hasher

	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return &p, salt, hash, nil
}

func (r *hasher) HashPassword(password string) (string, error) {
	salt, err := r.RandomBytes()
	if err != nil {
		return "", err
	}

	hashed := argon2.IDKey([]byte(password), salt, r.iterations, r.memory, r.parallelism, r.keyLength)

	return r.EncodeHash(hashed, salt), nil
}

func (r *hasher) CompareHash(hashedPassword, password string) (bool, error) {
	p, salt, hash, err := r.DecodeHash(hashedPassword)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}
