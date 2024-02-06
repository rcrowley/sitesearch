package main

import (
	"errors"

	"github.com/aws/smithy-go"
)

func awsErrorCode(err error) string {
	var ae smithy.APIError
	if errors.As(err, &ae) {
		return ae.ErrorCode()
	}
	return ""
}

func awsErrorCodeIs(err error, code string) bool {
	return awsErrorCode(err) == code
}
