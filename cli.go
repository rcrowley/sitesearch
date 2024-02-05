//go:build !lambda

package main

import (
	_ "embed"
	"log"
)

//go:generate env GOARCH=arm64 GOOS=linux go build -o bootstrap -tags lambda
//go:generate touch -t 202402040000.00 bootstrap
//go:generate zip -X lambda.zip bootstrap
//go:embed lambda.zip
var bin []byte

func main() {
	log.Print("main")
}
