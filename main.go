package main

import (
	"client/cmd"
	"client/helpers"
)

func init() {
	helpers.AddCommand(cmd.Help)
	helpers.AddCommand(cmd.Hello)
}

func main() {
	err := helpers.Run()
	if err != nil {
		println(err.Error())
	}
}
