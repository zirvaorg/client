package cmd

import (
	"client/helpers"
	"client/internal"
	"fmt"
)

var Login = &helpers.Command{
	Use:         "login",
	Description: "Login to the server",
	Run: func() {
		if helpers.CheckAuth() == true {
			fmt.Println("Already logged in")
			return
		}

		UUID := internal.CreateUUID()
		err := internal.WriteUUID(UUID)
		if err != nil {
			fmt.Println("Error writing UUID:", err)

		}

		println("UUID:", UUID)
	},
}
