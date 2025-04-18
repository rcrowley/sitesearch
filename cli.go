//go:build !lambda

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/rcrowley/mergician/files"
	"github.com/rcrowley/mergician/html"
	"github.com/rcrowley/sitesearch/index"
)

const (
	handler       = "bootstrap" // meet the provided.al* API
	runtime       = types.RuntimeProvidedal2023
	timeout int32 = 29 // seconds
)

func Main(args []string, stdin io.Reader, stdout io.Writer) {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	layout := flags.String("l", "", "site layout HTML document for search result pages")
	name := flags.String("n", "sitesearch", "name of the the Lambda function")
	region := flags.String("r", "", "AWS region to host the Lambda function")
	exclude := files.NewStringSliceFlag(flags, "x", "subdirectory of <docroot> to exclude (may be repeated)")
	flags.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: sitesearch -l <layout> [-n <name>] [-r <region>] [-x <exclude>[...]] [<docroot>[...]]
  -l <layout>   site layout HTML document for search result pages
  -n <name>     name of the the Lambda function (defaults to "sitesearch")
  -r <region>   AWS region to host the Lambda function (defaults to AWS_REGION or AWS_DEFAULT_REGION in the environment)
  -x <exclude>  subdirectory of <docroot> to exclude (may be repeated)
  <docroot>     document root directory to scan (defaults to the current working directory; may be repeated)

Synopsis: sitesearch scans each <docroot> (or the current working directory) for HTML documents, constructs an inverted index, and serves queries to it, including an HTML form and SERP, via AWS Lambda.
`)
	}
	flags.Parse(args[1:])
	if *layout == "" {
		log.Fatal("-l <layout> is required")
	}

	tmp, err := os.MkdirTemp("", "sitesearch-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	// Index all the HTML we've been told to. Store the index where the Lambda
	// function is eventually going to look for it.
	log.Printf("indexing HTML documents")
	f := func(n *html.Node) (string, string) {
		title := html.Text(html.Title(n)).String()

		// Remove redundant text from titles on from SERPs. TODO parameterize.
		title = strings.SplitN(title, "&mdash;", 2)[0]
		title = strings.SplitN(title, "—", 2)[0] // an unencoded &mdash;

		return strings.TrimSpace(title), strings.TrimSpace(html.Text(html.FirstParagraph(n)).String())
	}
	idx := must2(index.Open(filepath.Join(tmp, IdxFilename)))
	var docroots []string
	if flags.NArg() == 0 {
		docroots = []string{"."}
	} else {
		docroots = flags.Args()
	}
	lists := must2(files.AllInputs(docroots, *exclude))
	for _, list := range lists {
		must(idx.IndexHTMLFiles(list.QualifiedPaths(), f))
	}
	must(idx.Close())

	// Copy the search engine result page layout to where the Lambda function
	// is eventually going to look for it.
	must(os.WriteFile(
		filepath.Join(tmp, TmplFilename),
		must2(os.ReadFile(*layout)),
		0666,
	))

	// Package up the application, index, and layout for service in Lambda.
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

func init() {
	log.SetFlags(0)
}

func main() {
	Main(os.Args, os.Stdin, os.Stdout)
}
