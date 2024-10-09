package core

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
)

type UID struct {
	localID uint32
	role    uint32
}

func NewUID(localID uint32, role uint32) UID {
	return UID{
		localID: localID,
		role:    role,
	}
}

func (uid UID) String() string {
	value := uint64(uid.localID)<<31 | uint64(uid.role)<<30
	return base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", value)))
}

func (uid UID) GetLocalID() uint32 {
	return uid.localID
}

func (uid UID) GetRole() uint32 {
	return uid.role
}

func DecomposeUID(s string) (UID, error) {
	decodedStr, err := base64.RawStdEncoding.Strict().DecodeString(s)
	if err != nil {
		return UID{}, err
	}

	uid, err := strconv.ParseUint(string(decodedStr), 10, 64)
	if err != nil {
		return UID{}, err
	}

	if (1 << 31) > uid {
		return UID{}, errors.New("wrong uid")
	}

	u := UID{
		localID: uint32(uid >> 31),
		role:    uint32(uid >> 30 & 0x1),
	}

	return u, nil
}
