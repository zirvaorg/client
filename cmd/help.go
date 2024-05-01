package cmd

import "client/helpers"

var Help = &helpers.Command{
	Use:         "help",
	Description: "Prints help",
	Run: func() {
		for _, cmd := range helpers.Commands {
			println(cmd.Use, "-", cmd.Description)
		}
	},
}
