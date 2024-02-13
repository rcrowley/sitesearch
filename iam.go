package main

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func iamRole(ctx context.Context, cfg aws.Config, name string) (string, error) {
	client := iam.NewFromConfig(cfg)
	out, err := client.CreateRole(ctx, &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(`{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Action": "sts:AssumeRole",
			"Effect": "Allow",
			"Principal": {"Service": ["lambda.amazonaws.com"]}
		}
	]
}`),
		RoleName: aws.String(name),
		Tags: []types.Tag{{
			Key:   aws.String("Manager"),
			Value: aws.String("sitesearch"),
		}},
	})
	if awsErrorCodeIs(err, "EntityAlreadyExists") {
		out, err := client.GetRole(ctx, &iam.GetRoleInput{
			RoleName: aws.String(name),
		})
		if err != nil {
			return "", err
		}
		return aws.ToString(out.Role.Arn), nil
	} else if err != nil {
		return "", err
	}
	must(iam.NewRoleExistsWaiter(client).Wait(ctx, &iam.GetRoleInput{
		RoleName: aws.String(name),
	}, time.Minute))
	return aws.ToString(out.Role.Arn), nil
}
