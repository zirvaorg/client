package helpers

import (
	"archive/tar"
	"client/internal"
	"client/internal/package_url"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
)

type updateHelper struct {
	decompressor internal.Decompressor
}

var UpdateHelpers = &updateHelper{
	decompressor: internal.NewUnzip(),
}

func (u *updateHelper) ReplaceNewPackage(url string) error {
	tempDir := os.TempDir()
	tempZipFile := path.Join(tempDir, fmt.Sprintf("%s.%s", "zirva-client", package_url.ZIP_FILE_EXTENSION))

	err := u.downloadNewPackage(url, tempZipFile)
	if err != nil {
		return err
	}

	compressedFile, err := os.Open(tempZipFile)
	if err != nil {
		return err
	}
	defer compressedFile.Close()

	createdFileName, err := internal.FilenameWithExtension("client")

	if err != nil {
		return err
	}

	createdFile, err := os.Create(path.Join(tempDir, createdFileName))

	if err != nil {
		return err
	}

	if err = os.Chmod(createdFile.Name(), 0777); err != nil {
		return err
	}
	defer createdFile.Close()

	fmt.Println("Package has been downloaded. Decompressing now.")
	if err = u.decompressor.Decompress(createdFile, compressedFile); err != nil {
		return err
	}
	fmt.Println("Decompress complete.")

	currentApp, err := os.Executable()

	if err != nil {
		return err
	}

	fmt.Println("Renaming...")
	if err = os.Rename(createdFile.Name(), currentApp); err != nil {
		return err
	}
	fmt.Println("Setting chmod")
	if err = os.Chmod(currentApp, 0755); err != nil {
		return err
	}
	fmt.Printf("Main process PID: %d\n", os.Getpid())

	fmt.Println("Running new package.")
	if err = u.runNewPackage(currentApp, "--stealth"); err != nil {
		return err
	}

	return nil
}

func (u *updateHelper) unzip(tempZipFile io.Reader) error {
	r := tar.NewReader(tempZipFile)
	log.Fatal(r.Next())
	return nil
}

func (u *updateHelper) downloadNewPackage(url, filePath string) error {
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

func (u *updateHelper) runNewPackage(params ...string) error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	args := append([]string{executable})
	fmt.Println(args)
	cmd := exec.Command("open", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err = cmd.Start()

	if err != nil {
		return err
	}
	fmt.Printf("New PID: %d\n", cmd.Process.Pid)

	return nil
}

func (u *updateHelper) ReleaseCurrentProcess() error {
	pid := os.Getpid()
	fmt.Printf("Current Process: %d\n", pid)
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	fmt.Printf("Releasing process %d\n", pid)
	err = p.Release()
	fmt.Printf("Released process %d\n", pid)
	if err != nil {
		return err
	}
	return nil
}

func (u *updateHelper) KillCurrentProcess() error {
	pid := os.Getpid()
	fmt.Printf("Current Process: %d\n", pid)
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	fmt.Printf("Killing process %d\n", pid)
	err = p.Kill()
	fmt.Printf("Killed process %d\n", pid)
	if err != nil {
		return err
	}
	return nil
}
