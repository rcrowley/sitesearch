//go:build !lambda

package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
)

// Zip writes the starting-point zip file that contains the compiled Lambda
// function, named bootstrap per the Lambda provided.al2 runtime requirement,
// to disk, and adds the file with the given relative pathname to it.
func Zip(pathname string) (string, error) {
	const zipFilename = "sitesearch.zip"

	f, err := os.Create(zipFilename)
	if err != nil {
		return "", err
	}
	if _, err := f.Write(bootstrapZip); err != nil {
		if err2 := f.Close(); err2 != nil {
			return "", fmt.Errorf("%w %w", err, err2)
		}
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}

	cmd := exec.Command("zip", "-X", "-r", zipFilename, pathname)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return zipFilename, nil
}

//go:generate env GOARCH=arm64 GOOS=linux go build -o bootstrap -tags lambda
//go:generate touch -t 202402040000.00 bootstrap
//go:generate zip -X bootstrap.zip bootstrap
//go:embed bootstrap.zip
var bootstrapZip []byte
