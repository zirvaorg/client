package cmd

import (
	"client/helpers"
	"client/internal/package_url"
	"fmt"
)

const (
	nextVersion = "v0.0.9"
)

var Update = &helpers.Command{
	Use:         "update",
	Description: "Update the application",
	Run: func() {
		fmt.Println("Updating...")
		packageURL := fmt.Sprintf(package_url.PackageURLFormat, nextVersion, nextVersion)
		err := helpers.UpdateHelpers.ReplaceNewPackage(packageURL)
		if err != nil {
			fmt.Println("Update failed:", err)
			return
		}

		//sigs := make(chan os.Signal, 1)
		//signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		//
		//fmt.Println("Update successful. Restarting the application...")
		//go func() {
		//	<-sigs
		//	err := helpers.UpdateHelpers.KillCurrentProcess()
		//	if err != nil {
		//		fmt.Println("Old process couldn't be kill:", err)
		//	}
		//}()

		err = helpers.UpdateHelpers.KillCurrentProcess()

		if err != nil {
			fmt.Println("Main process release failed:", err)
		}

		fmt.Println("Update succeeded, main process has been released")
		select {}
	},
}
