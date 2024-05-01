package internal

import (
	"encoding/base64"
	"github.com/google/uuid"
	"os"
)

const UUIDFile = "uuid.conf"

func CreateUUID() string {
	uuId := uuid.New()
	uuIdBytes := uuId[:]
	return base64.StdEncoding.EncodeToString(append(uuIdBytes))
}

func WriteUUID(UUID string) error {
	err := os.WriteFile(UUIDFile, []byte(UUID), 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadUUID() (string, error) {
	data, err := os.ReadFile(UUIDFile)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func DeleteUUID() error {
	err := os.Remove(UUIDFile)
	if err != nil {
		return err
	}
	return nil
}
