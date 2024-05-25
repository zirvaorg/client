package helpers

import (
	"client/internal"
	"client/internal/observer"
	"client/internal/package_url"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-version"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
)

type release struct {
	TagName string `json:"tag_name"`
}

type UpdateHelper struct {
	decompressor     internal.Decompressor
	checksumVerifier *internal.Checksum
	observer         observer.Observer
}

const (
	githubReleaseApiUri = "https://api.github.com/repos/zirvaorg/client/releases/latest"
	tempZipFileName     = "zirva-client"
)

var (
	UpdateHelpers = &UpdateHelper{
		decompressor:     internal.NewUnzip(),
		checksumVerifier: internal.NewChecksum(sha256.New()),
		observer:         observer.Console{},
	}
	LatestVersion          *version.Version
	ErrLatestVersionError  = errors.New("something went wrong when trying to determine the latest version")
	ErrChecksumVerifyError = errors.New("failed to verify checksum")
)

func init() {
	latestRelease, err := getLatestVersion()
	if err != nil {
		log.Fatal(ErrLatestVersionError)
		return
	}
	v, err := version.NewVersion(latestRelease.TagName)
	if err != nil {
		log.Fatal(ErrLatestVersionError)
		return
	}
	LatestVersion = v
}

func getLatestVersion() (*release, error) {
	resp, err := http.Get(githubReleaseApiUri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var releaseResp *release
	if err := json.NewDecoder(resp.Body).Decode(&releaseResp); err != nil {
		return nil, err
	}

	return releaseResp, nil
}

func (u *UpdateHelper) IsUpToDate(currentVersion string, latestVersion *version.Version) (bool, error) {
	current, err := version.NewVersion(currentVersion)
	if err != nil {
		return false, err
	}

	return current.GreaterThanOrEqual(latestVersion), nil // normally we should use just Equal but this way is safer.
}

func (u *UpdateHelper) ReplaceNewPackage(packageUrl, checksumUrl string) error {
	tempDir := os.TempDir()
	tempZipFile := path.Join(tempDir, fmt.Sprintf("%s.%s", tempZipFileName, package_url.ZIP_FILE_EXTENSION))

	err := u.downloadNewPackage(packageUrl, tempZipFile)
	if err != nil {
		return err
	}

	if isChecksumVerified, err := u.checksumVerifier.Verify(checksumUrl, tempZipFile); !isChecksumVerified {
		_ = os.Remove(tempZipFile)
		if err != nil {
			return err
		}
		return ErrChecksumVerifyError
	}

	u.notifyString("Checksum has been verified!")

	createdFileName := internal.FilenameWithExtension("client")
	createdFile, err := os.Create(path.Join(tempDir, createdFileName))

	if err != nil {
		return err
	}

	if err = os.Chmod(createdFile.Name(), 0555); err != nil {
		return err
	}
	defer createdFile.Close()

	u.notifyString("Package has been downloaded. Decompressing now.")
	if err = u.decompressor.Decompress(createdFile, tempZipFile); err != nil {
		return err
	}
	u.notifyString("Decompress complete.")

	currentApp, err := os.Executable()

	if err != nil {
		return err
	}

	u.notifyString("Renaming...")
	if err = os.Rename(createdFile.Name(), currentApp); err != nil {
		return err
	}
	u.notifyString("Setting chmod")
	if err = os.Chmod(currentApp, 0755); err != nil {
		return err
	}

	u.notifyString("Running new package.")
	if err = u.runNewPackage(); err != nil {
		return err
	}

	return nil
}

func (u *UpdateHelper) notifyString(msg string) {
	if u.observer != nil {
		u.observer.Update(msg)
	}
}

func (u *UpdateHelper) downloadNewPackage(url, filePath string) error {
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

func (u *UpdateHelper) runNewPackage() error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	// we can't use --stealth param since "open" doesn't allow.
	//TODO:: find a solution if possible.
	cmd := exec.Command("open", executable)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err = cmd.Start()

	if err != nil {
		return err
	}

	return nil
}

func (u *UpdateHelper) KillCurrentProcess() error {
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
