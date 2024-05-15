package cmd

import (
	"client/helpers"
	"fmt"
)

var Start = &helpers.Command{
	Use:         "start",
	Description: "Start the client",
	Run: func() {
		if !helpers.CheckAuth() {
			fmt.Println("Not logged in")
			return
		}

		fmt.Println("Starting client")
	},
}
