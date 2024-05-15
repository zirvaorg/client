package main

import (
	"client/cmd"
	"client/helpers"
	"flag"
	"fmt"
	"github.com/fatih/color"
)

var stealth *bool

func init() {
	helpers.AddCommand(cmd.Help)
	helpers.AddCommand(cmd.Start)
	helpers.AddCommand(cmd.Login)
	helpers.AddCommand(cmd.Logout)
	helpers.AddCommand(cmd.Update)
	helpers.AddCommand(cmd.Exit)

	stealth = flag.Bool("stealth", false, "run in stealth mode") // stealth mode for updating
	flag.Parse()
}

func main() {
	if *stealth {
		if helpers.CheckAuth() {
			cmd.Start.Run()
			return
		}
	}

	b := color.New(color.FgHiYellow)
	b.Println("┌──────────────────────────────────────┐")
	b.Println("│ ███████╗██╗██████╗ ██╗   ██╗ █████╗  │")
	b.Println("│ ╚══███╔╝██║██╔══██╗██║   ██║██╔══██╗ │")
	b.Println("│   ███╔╝ ██║██████╔╝██║   ██║███████║ │")
	b.Println("│  ███╔╝  ██║██╔══██╗╚██╗ ██╔╝██╔══██║ │")
	b.Println("│ ███████╗██║██████╔╝ ╚████╔╝ ██║  ██║ │")
	b.Println("│ ╚══════╝╚═╝╚═════╝   ╚═══╝  ╚═╝  ╚═╝ │")
	b.Println("└──────────────────────────────────────┘")
	fmt.Print("zirva client v0.1.0 (https://zirva.org)\n")
	info := color.New(color.FgHiYellow).SprintFunc()
	fmt.Printf("type '%s' for a list of commands\n", info("help"))

	err := helpers.Run()
	if err != nil {
		println(err.Error())
	}
}
