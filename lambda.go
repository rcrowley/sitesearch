//go:build lambda

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {

	// Move to the writable part of the filesystem before handling requests.
	must(CopyRecursive(IdxFilename, filepath.Join("/tmp", IdxFilename)))
	must(CopyRecursive(TmplFilename, filepath.Join("/tmp", TmplFilename)))
	must(os.Chdir("/tmp"))

	lambda.Start(SearchHandler)
}
