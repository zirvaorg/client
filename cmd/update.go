package cmd

import (
	"client/helpers"
	"client/internal/package_url"
	"fmt"
)

const (
	currentVersion              = "v0.0.8"
	latestChecksumFileUriFormat = "https://github.com/zirvaorg/client/releases/download/%s/client_%s_checksums.txt"
)

var Update = &helpers.Command{
	Use:         "update",
	Description: "Update the application",
	Run: func() {
		fmt.Println("Current Version:", currentVersion)
		if helpers.LatestVersion.Original() == "" {
			fmt.Println("Error getting latest version")
			return
		}
		fmt.Printf("Latest version: %s\n", helpers.LatestVersion)

		isUpToDate, err := helpers.UpdateHelpers.IsUpToDate(currentVersion, helpers.LatestVersion)
		if err != nil {
			fmt.Println("Error comparing two versions")
			return
		}

		if isUpToDate {
			fmt.Printf("Application is up to date, current: %s, latest: %s\n", currentVersion, helpers.LatestVersion.Original())
			return
		}

		fmt.Println("Updating...")

		packageURL := fmt.Sprintf(package_url.PackageURLFormat, helpers.LatestVersion.Original(), helpers.LatestVersion.Original())
		checksumUrl := fmt.Sprintf(latestChecksumFileUriFormat, helpers.LatestVersion.Original(), helpers.LatestVersion.String())

		err = helpers.UpdateHelpers.ReplaceNewPackage(packageURL, checksumUrl)
		if err != nil {
			fmt.Println("Update failed:", err)
			return
		}

		err = helpers.UpdateHelpers.KillCurrentProcess()

		if err != nil {
			fmt.Println("Main process release failed:", err)
		}

		fmt.Println("Update succeeded, main process has been released")
		select {}
	},
}
