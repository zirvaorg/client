package cmd

import (
	"client/helpers"
	"fmt"
	"github.com/fatih/color"
)

var Brand = &helpers.Command{
	Run: func() {
		fmt.Print("#  ███████╗██╗██████╗ ██╗   ██╗ █████╗\n")
		fmt.Print("#  ╚══███╔╝██║██╔══██╗██║   ██║██╔══██╗\n")
		fmt.Print("#    ███╔╝ ██║██████╔╝██║   ██║███████║\n")
		fmt.Print("#   ███╔╝  ██║██╔══██╗██║   ██║██╔══██║\n")
		fmt.Print("#  ███████╗██║██║  ██║╚██████╔╝██║  ██║\n")
		fmt.Print("#  ╚══════╝╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝\n")
		fmt.Print("#  zirva client v0.1.0 (https://zirva.org)\n")
		info := color.New(color.FgHiYellow).SprintFunc()
		fmt.Printf("#  type '%s' for a list of commands\n", info("help"))
	},
}
