package main

import (
	"client/cmd"
	"client/helpers"
)

func init() {
	helpers.AddCommand(cmd.Help)
	helpers.AddCommand(cmd.Hello)
	helpers.AddCommand(cmd.Login)
}

func main() {
	cmd.Brand.Run()

	err := helpers.Run()
	if err != nil {
		println(err.Error())
	}
}
