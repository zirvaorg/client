package internal

import (
	"encoding/base64"
	"github.com/google/uuid"
)

func CreateUUID() string {
	uuId := uuid.New()
	uuIdBytes := uuId[:]
	return base64.StdEncoding.EncodeToString(append(uuIdBytes))
}
