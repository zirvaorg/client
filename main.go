package main

import (
	"client/cmd"
	"client/helpers"
	"flag"
)

func init() {
	helpers.AddCommand(cmd.Help)
	helpers.AddCommand(cmd.Login)
	helpers.AddCommand(cmd.Logout)
	helpers.AddCommand(cmd.Update)
	helpers.AddCommand(cmd.Exit)

	hideBrand := flag.Bool("hide-brand", false, "hide the brand")
	flag.Parse()

	if !*hideBrand {
		cmd.Brand.Run()
	}
}

func main() {
	err := helpers.Run()
	if err != nil {
		println(err.Error())
	}
}
