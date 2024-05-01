package cmd

import (
	"client/helpers"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// PackageURL @todo: change the package URL to github
const (
	PackageURL = "https://sercanarga.com/client"
)

var Update = &helpers.Command{
	Use:         "update",
	Description: "Update the application",
	Run: func() {
		u := &helpers.UpdateHelpers{}

		err := u.ReplaceNewPackage(PackageURL)
		if err != nil {
			fmt.Println("Update failed:", err)
			return
		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigs
			err := u.KillCurrentProcess()
			if err != nil {
				fmt.Println("Old process couldn't be kill:", err)
			}
		}()

		select {}
	},
}
