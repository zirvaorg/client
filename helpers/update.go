package helpers

import (
	"io"
	"net/http"
	"os"
	"os/exec"
)

type UpdateHelpers struct{}

func (u *UpdateHelpers) ReplaceNewPackage(url string) error {
	tempFile := "/tmp/zirva-client"

	err := u.downloadNewPackage(url, tempFile)
	if err != nil {
		return err
	}

	currentApp := os.Args[0]

	err = os.Rename(tempFile, currentApp)
	if err != nil {
		return err
	}

	err = os.Chmod(currentApp, 0755)
	if err != nil {
		return err
	}

	err = u.runNewPackage()
	if err != nil {
		return err
	}

	return nil
}

func (u *UpdateHelpers) downloadNewPackage(url, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func (u *UpdateHelpers) runNewPackage() error {
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (u *UpdateHelpers) KillCurrentProcess() error {
	pid := os.Getpid()
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = p.Kill()
	if err != nil {
		return err
	}
	return nil
}
