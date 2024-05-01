package cmd

import (
	"client/helpers"
	"client/internal"
)

var Hello = &helpers.Command{
	Use:         "hello",
	Description: "Prints hello",
	Run: func() {
		os, _ := internal.DetectOS()
		println("Hello from", os)
	},
}
