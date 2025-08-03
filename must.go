package main

import (
	"log"
	"os"
)

func must(err error) {
	if err != nil {
		log.Output(2, err.Error())
		os.Exit(1)
	}
}

func must2[T any](v T, err error) T {
	if err != nil {
		log.Output(2, err.Error())
		os.Exit(1)
	}
	return v
}
