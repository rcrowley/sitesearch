//go:build !lambda

package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
)

const ZipFilename = "sitesearch.zip"

// Zip writes the starting-point zip file that contains the compiled Lambda
// function, named bootstrap per the Lambda provided.al* runtime requirement,
// to disk, and adds the file with the given relative pathname to it.
func Zip(
	zipPath string, // output
	idxPath, tmplPath string, // input
) ([]byte, error) {

	f, err := os.Create(zipPath)
	if err != nil {
		return nil, err
	}
	if _, err := f.Write(bootstrapZip); err != nil {
		if err2 := f.Close(); err2 != nil {
			return nil, fmt.Errorf("%w %w", err, err2)
		}
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}

	// If you ever need to debug zip(1), add these options:
	// "-la", "-lf", "/tmp/zip.log", "-li"

	if err := exec.Command("zip", "-X", "-r", zipPath, idxPath).Run(); err != nil {
		return nil, err
	}

	if err := exec.Command("zip", "-X", "-j", zipPath, tmplPath).Run(); err != nil {
		return nil, err
	}

	return os.ReadFile(zipPath)
}

//go:generate env GOARCH=arm64 GOOS=linux go build -o bootstrap -tags lambda
//go:generate touch -t 202402040000.00 bootstrap
//go:generate zip -X bootstrap.zip bootstrap
//go:embed bootstrap.zip
var bootstrapZip []byte
