package cmd

import (
	"client/helpers"
	"client/internal"
	"fmt"
)

var Logout = &helpers.Command{
	Use:         "logout",
	Description: "Logout to the server",
	Run: func() {
		if helpers.CheckAuth() != true {
			fmt.Println("Not logged in")
			return
		}

		err := internal.DeleteUUID()
		if err != nil {
			fmt.Println("Error deleting UUID:", err)
		}

		println("Logged out")
	},
}
