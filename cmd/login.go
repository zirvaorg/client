package cmd

import (
	"client/helpers"
	"client/internal"
)

var Login = &helpers.Command{
	Use:         "login",
	Description: "Login to the server",
	Run: func() {
		UUID := internal.CreateUUID()
		println("UUID:", UUID)
	},
}
