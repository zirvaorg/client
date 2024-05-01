package main

import (
	"client/cmd"
	"client/helpers"
)

func init() {
	helpers.AddCommand(cmd.Help)
	helpers.AddCommand(cmd.Login)
	helpers.AddCommand(cmd.Logout)
	helpers.AddCommand(cmd.Update)
}

func main() {
	cmd.Brand.Run()

	err := helpers.Run()
	if err != nil {
		println(err.Error())
	}
}
