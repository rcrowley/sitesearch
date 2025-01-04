package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestLambdaForm(t *testing.T) {
	ctx := context.Background()
	req := events.LambdaFunctionURLRequest{RawPath: "/search/"}
	resp, err := SearchHandler(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	actual := resp.Body
	expected := `<!DOCTYPE html>
<html>
<head>
<title>Sitesearch</title>
</head>
<body>
<div class="body">
    <form class="sitesearch">
        <input name="q" placeholder="query" type="text" value=""/>
        <input type="submit" value="Search"/>
    </form>
</div>
</body>
</html>
`
	if actual != expected {
		t.Fatalf("actual: %s != expected: %s", actual, expected)
	}
}

// TestLambdaSERP tests SearchHandler with a ?q=cool querystring such that it
// tries to execute a search. It doesn't work because there's no index but that
// doesn't matter because the actual search is tested elsewhere.
func TestLambdaSERP(t *testing.T) {
	ctx := context.Background()
	req := events.LambdaFunctionURLRequest{
		RawPath:               "/search/",
		QueryStringParameters: map[string]string{"q": "cool"},
	}
	resp, err := SearchHandler(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	actual := resp.Body
	expected := `<!DOCTYPE html>
<html>
<head>
<title>Sitesearch</title>
</head>
<body>
<div class="body">
    <pre class="sitesearch">cannot open index, path does not exist</pre>
</div>
</body>
</html>
`
	if actual != expected {
		t.Fatalf("actual: %s != expected: %s", actual, expected)
	}
}
