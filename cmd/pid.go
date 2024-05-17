package cmd

import (
	"client/helpers"
	"fmt"
	"github.com/fatih/color"
	"os"
)

var Pid = &helpers.Command{
	Use:         "pid",
	Description: "Prints current PID",
	Run: func() {
		key := color.New(color.FgHiYellow).SprintFunc()
		desc := color.New(color.FgWhite).SprintFunc()

		fmt.Printf("%s	%s\n", key("Current PID:"), desc(os.Getpid()))
	},
}
