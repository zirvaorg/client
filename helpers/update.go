package helpers

import (
	"bufio"
	"client/internal"
	"client/internal/package_url"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-version"
	"hash"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

type release struct {
	TagName string `json:"tag_name"`
}

const (
	githubReleaseApiUri         = "https://api.github.com/repos/zirvaorg/client/releases/latest"
	bufferSize                  = 65536
	latestChecksumFileUriFormat = "https://github.com/zirvaorg/client/releases/download/%s/client_%s_checksums.txt"
)

var (
	UpdateHelpers = &updateHelper{
		decompressor: internal.NewUnzip(),
	}
	LatestVersion *version.Version
)

func init() {
	latestRelease, err := getLatestVersion()
	if err != nil {
		log.Fatal("Something went wrong when trying to determine the latest version")
		return
	}
	v, err := version.NewVersion(latestRelease.TagName)
	if err != nil {
		log.Fatal("Something went wrong when trying to determine the version")
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

type updateHelper struct {
	decompressor internal.Decompressor
}

func (u *updateHelper) IsUpToDate(currentVersion string, latest *version.Version) (bool, error) {
	current, err := version.NewVersion(currentVersion)
	if err != nil {
		return false, err
	}

	return current.GreaterThanOrEqual(latest), nil // normally we should use just Equal but this way is safer.
}

func (u *updateHelper) calculateChecksum(hashAlgorithm hash.Hash, reader io.Reader) (string, error) {
	buf := make([]byte, bufferSize)

	for {
		switch n, err := reader.Read(buf); err {
		case nil:
			hashAlgorithm.Write(buf[:n])
		case io.EOF:
			return fmt.Sprintf("%x", hashAlgorithm.Sum(nil)), nil
		default:
			return "", errors.New("failed to calculate checksum")
		}
	}
}

func (u *updateHelper) compareChecksum(tempZipFilePath string) (bool, error) {
	// create new file pointer so reading process won't affect the other one
	reader, err := os.Open(tempZipFilePath)
	if err != nil {
		return false, err
	}
	actualChecksum, err := u.calculateChecksum(sha256.New(), reader)
	if err != nil {
		return false, err
	}
	expectedChecksum, err := u.getChecksumFromGithub()
	if err != nil {
		return false, err
	}

	fmt.Println(actualChecksum, expectedChecksum, actualChecksum == expectedChecksum)

	return actualChecksum == expectedChecksum, nil
}

func (u *updateHelper) getChecksumFromGithub() (string, error) {
	checksumUrl := fmt.Sprintf(latestChecksumFileUriFormat, LatestVersion.Original(), LatestVersion.String())

	resp, err := http.Get(checksumUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fileScanner := bufio.NewScanner(resp.Body)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := strings.Split(fileScanner.Text(), "  ")
		if len(line) != 2 || !strings.HasPrefix(line[1], "zirva") {
			return "", errors.New("failed to parse checksum file")
		}
		if line[1] == fmt.Sprintf(package_url.ZipFileNameFormat, LatestVersion.Original()) {
			fmt.Println("Related checksum was found from github")
			return line[0], nil
		}
	}
	return "", errors.New("failed to find checksum")
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
	defer func(compressedFile *os.File) {
		_ = compressedFile.Close()
	}(compressedFile)

	if isChecksumVerified, err := u.compareChecksum(tempZipFile); !isChecksumVerified {
		fmt.Println("Checksum verification failed, downloaded file is deleting now.")
		_ = compressedFile.Close()
		_ = os.Remove(tempZipFile)
		if err != nil {
			return err
		}
		return errors.New("checksum verification failed")
	}

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

	fmt.Println("Running new package.")
	if err = u.runNewPackage(); err != nil {
		return err
	}

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

func (u *updateHelper) runNewPackage() error {
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

func (u *updateHelper) KillCurrentProcess() error {
	pid := os.Getpid()
	fmt.Printf("Current Process: %d\n", pid)
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
