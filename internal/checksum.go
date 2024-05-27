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
	bufferSize = 65536
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

func (c *Checksum) Verify(checksumUrl, filePath string, latestVersion version.Version) (bool, error) {
	return c.compareChecksum(checksumUrl, filePath, latestVersion)
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

func (c *Checksum) compareChecksum(checksumUrl, filePath string, latestVersion version.Version) (bool, error) {
	// create new file pointer so reading process won't affect the other one
	reader, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	actualChecksum, err := c.calculateChecksum(reader)
	if err != nil {
		return false, err
	}
	expectedChecksum, err := c.getChecksumFromGithub(checksumUrl, latestVersion)
	if err != nil {
		return false, err
	}

	return actualChecksum == expectedChecksum, nil
}

func (c *Checksum) getChecksumFromGithub(checksumUrl string, latestVersion version.Version) (string, error) {
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
			return line[0], nil
		}
	}
	return "", ErrInvalidChecksumError
}
