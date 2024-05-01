package cmd

import (
	"client/helpers"
	"fmt"
	"github.com/fatih/color"
)

var Help = &helpers.Command{
	Use:         "help",
	Description: "Prints help",
	Run: func() {
		key := color.New(color.FgHiYellow).SprintFunc()
		desc := color.New(color.FgWhite).SprintFunc()

		for _, cmd := range helpers.Commands {
			if cmd.Use == "help" {
				continue
			}

			fmt.Printf("	%s		%s\n", key(cmd.Use), desc(cmd.Description))
		}
	},
}
