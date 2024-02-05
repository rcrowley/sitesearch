//go:build lambda

package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	lambda.Start(SearchHandler)
}
