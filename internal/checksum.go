package internal

import (
	"bufio"
	"client/internal/package_url"
	"errors"
	"fmt"
	"github.com/hashicorp/go-version"
	"hash"
	"io"
	"net/http"
	"os"
	"strings"
)

type Checksum struct {
	HashAlgorithm hash.Hash
}

const (
	bufferSize                  = 65536
	latestChecksumFileUriFormat = "https://github.com/zirvaorg/client/releases/download/%s/client_%s_checksums.txt"
)

var (
	ErrCalculateChecksumError = errors.New("failed to calculate checksum")
	ErrParseChecksumFileError = errors.New("failed to parse checksum file")
	ErrInvalidChecksumError   = errors.New("failed to find checksum")
)

func NewChecksum(hashAlgorithm hash.Hash) *Checksum {
	return &Checksum{
		HashAlgorithm: hashAlgorithm,
	}
}

func (c *Checksum) Verify(latestVersion *version.Version, filePath string) (bool, error) {
	return c.compareChecksum(latestVersion, filePath)
}

func (c *Checksum) calculateChecksum(reader io.Reader) (string, error) {
	defer c.HashAlgorithm.Reset()
	buf := make([]byte, bufferSize)

	for {
		switch n, err := reader.Read(buf); err {
		case nil:
			c.HashAlgorithm.Write(buf[:n])
		case io.EOF:
			return fmt.Sprintf("%x", c.HashAlgorithm.Sum(nil)), nil
		default:
			return "", ErrCalculateChecksumError
		}
	}
}

func (c *Checksum) compareChecksum(latestVersion *version.Version, filePath string) (bool, error) {
	// create new file pointer so reading process won't affect the other one
	reader, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	actualChecksum, err := c.calculateChecksum(reader)
	if err != nil {
		return false, err
	}
	expectedChecksum, err := c.getChecksumFromGithub(latestVersion)
	if err != nil {
		return false, err
	}

	return actualChecksum == expectedChecksum, nil
}

func (c *Checksum) getChecksumFromGithub(latestVersion *version.Version) (string, error) {
	checksumUrl := fmt.Sprintf(latestChecksumFileUriFormat, latestVersion.Original(), latestVersion.String())

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
			return "", ErrParseChecksumFileError
		}
		if line[1] == fmt.Sprintf(package_url.ZipFileNameFormat, latestVersion.Original()) {
			fmt.Println("Related checksum was found from github")
			return line[0], nil
		}
	}
	return "", ErrInvalidChecksumError
}
