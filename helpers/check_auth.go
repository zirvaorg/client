package helpers

import "client/internal"

func CheckAuth() bool {
	check, _ := internal.ReadUUID()
	if check != "" {
		return true
	}

	return false
}
