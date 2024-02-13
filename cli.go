//go:build !lambda

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/rcrowley/sitesearch/index"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	handler       = "bootstrap" // meet the provided.al* API
	runtime       = types.RuntimeProvidedal2023
	timeout int32 = 29 // seconds
)

func init() {
	log.SetFlags(0)
}

func main() {
	name := flag.String("n", "sitesearch", "name of the the Lambda function")
	region := flag.String("r", "", "AWS region to host the Lambda function")
	tmpl := flag.String("t", "", "HTML template for search result pages")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: sitesearch [-n <name>] [-r <region>] -t <template> <input>[...]
  -n <name>      name of the the Lambda function (default "sitesearch")
  -r <region>    AWS region to host the Lambda function (default to AWS_DEFAULT_REGION in the environment)
  -t <template>  HTML template for search result pages
  <input>[...]   pathname, relative to your site's root, of one or more HTML files, given as command-line arguments or on standard input
`)
	}
	flag.Parse()
	if *tmpl == "" {
		log.Fatal("-t <template> is required")
	}

	tmp, err := os.MkdirTemp("", "sitesearch-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	// Index all the HTML we've been told to. Store the index where the Lambda
	// function is eventually going to look for it.
	log.Printf("indexing HTML documents")
	idx := must2(index.Open(filepath.Join(tmp, IdxFilename)))
	must(idx.IndexHTMLFiles(flag.Args()))
	if !terminal.IsTerminal(0) {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			must(idx.IndexHTMLFile(scanner.Text()))
		}
		must(scanner.Err())
	}
	must(idx.Close())

	// Copy the search engine result page template to where the Lambda function
	// is eventually going to look for it.
	must(os.WriteFile(
		filepath.Join(tmp, TmplFilename),
		must2(os.ReadFile(*tmpl)),
		0666,
	))

	// Package up the application, index, and template for service in Lambda.
	log.Printf("packaging the search application")
	oldpwd := must2(os.Getwd())
	must(os.Chdir(tmp))
	zipFile := must2(Zip(ZipFilename, IdxFilename, TmplFilename))
	must(os.Chdir(oldpwd))

	// Create or update a Lambda function to serve this search application.
	// Use whatever AWS credentials we find lying around and the region either
	// found in the environment or given as an option.
	ctx := context.Background()
	var options []func(*config.LoadOptions) error
	if *region != "" {
		options = append(options, config.WithRegion(*region))
	}
	cfg := must2(config.LoadDefaultConfig(ctx, options...))
	log.Printf("creating or updating the IAM role")
	roleARN := must2(iamRole(ctx, cfg, "sitesearch"))
	log.Printf("creating or updating the Lambda function")
	client := lambda.NewFromConfig(cfg)
	_, err = client.CreateFunction(ctx, &lambda.CreateFunctionInput{
		Architectures: []types.Architecture{types.ArchitectureArm64},
		Code:          &types.FunctionCode{ZipFile: zipFile},
		FunctionName:  name,
		Handler:       aws.String(handler),
		PackageType:   types.PackageTypeZip,
		Role:          aws.String(roleARN),
		Runtime:       runtime,
		Tags:          map[string]string{"Manager": "sitesearch"},
		Timeout:       aws.Int32(timeout),
	})
	if awsErrorCodeIs(err, "ResourceConflictException") {
		/*
			must2(client.UpdateFunctionConfiguration(
				ctx,
				&lambda.UpdateFunctionConfigurationInput{
					FunctionName: name,
					Handler:      aws.String(handler),
					Role:         aws.String(roleARN),
					Runtime:      runtime,
					Timeout:      aws.Int32(timeout),
				},
			))
		*/
		must2(client.UpdateFunctionCode(
			ctx,
			&lambda.UpdateFunctionCodeInput{
				Architectures: []types.Architecture{types.ArchitectureArm64},
				Publish:       true,
				ZipFile:       zipFile,
				FunctionName:  name,
			},
		))
	} else if err != nil {
		log.Fatal(err)
	}
	log.Printf("publishing the Lambda function URL")
	if _, err := client.AddPermission(
		ctx,
		&lambda.AddPermissionInput{
			Action:              aws.String("lambda:InvokeFunctionUrl"),
			FunctionName:        name,
			FunctionUrlAuthType: types.FunctionUrlAuthTypeNone,
			Principal:           aws.String("*"),
			StatementId:         aws.String("sitesearch"),
		},
	); err != nil && !awsErrorCodeIs(err, "ResourceConflictException") {
		log.Fatal(err)
	}
	out, err := client.CreateFunctionUrlConfig(
		ctx,
		&lambda.CreateFunctionUrlConfigInput{
			AuthType:     types.FunctionUrlAuthTypeNone,
			FunctionName: name,
		},
	)
	var functionURL string
	if err == nil {
		functionURL = aws.ToString(out.FunctionUrl)
	} else if awsErrorCodeIs(err, "ResourceConflictException") {
		out := must2(client.GetFunctionUrlConfig(
			ctx,
			&lambda.GetFunctionUrlConfigInput{
				FunctionName: name,
			},
		))
		functionURL = aws.ToString(out.FunctionUrl)
	} else {
		log.Fatal(err)
	}
	fmt.Println(functionURL)
}
