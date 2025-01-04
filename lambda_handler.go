package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rcrowley/mergician/html"
	"golang.org/x/net/html/atom"
)

func SearchHandler(ctx context.Context, req events.LambdaFunctionURLRequest) (resp events.LambdaFunctionURLResponse, err error) {
	resp.StatusCode = http.StatusOK
	resp.Headers = map[string]string{"Content-Type": "text/html; charset=utf-8"}
	var n, tmpl *html.Node

	if !strings.HasSuffix(req.RawPath, "/") {
		resp.StatusCode = http.StatusFound
		resp.Headers["Location"] = fmt.Sprintf("%s/", req.RawPath)
		resp.Body = html.String(errorNode(fmt.Errorf("redirecting to %s", req.RawPath)))
		return
	}

	if q := req.QueryStringParameters["q"]; q == "" {
		n = &html.Node{
			DataAtom: atom.Body,
			Type:     html.ElementNode,
		}
		n.AppendChild(newlineIndent(1))
		n.AppendChild(form(q))
		n.AppendChild(newlineIndent(0))
	} else {
		n, err = Search(q)
	}
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		n = errorNode(err)
	}

	if tmpl, err = html.ParseFile(TmplFilename); err == nil {
		if n, err = html.Merge([]*html.Node{tmpl, n}, html.DefaultRules()); err != nil {
			return
		}
	}

	resp.Body = html.String(n)
	return
}
