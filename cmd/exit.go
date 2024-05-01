package cmd

import "client/helpers"

var Exit = &helpers.Command{
	Use:         "exit",
	Description: "Exit the application",
	Run: func() {
		err := helpers.UpdateHelpers.KillCurrentProcess()
		if err != nil {
			println("Error while exiting:", err)
		}
	},
}
