Sitesearch
==========

Index all the HTML files in a document root directory and serve queries to it, including an HTML form and SERP, via AWS Lambda.

Installation
------------

```sh
git clone https://github.com/rcrowley/sitesearch.git
cd sitesearch
go generate
go install
```

Usage
-----

```sh
sitesearch -l <layout> [-n <name>] [-r <region>] [-x <exclude>[...]] [<docroot>[...]]
```

* `-l <layout>`: site layout HTML document for search result pages
* `-n <name>`: name of the the Lambda function (defaults to "sitesearch")
* `-r <region>`: AWS region to host the Lambda function (defaults to `AWS_REGION` or `AWS_DEFAULT_REGION` in the environment)
* `-x <exclude>`: subdirectory of `<docroot>` to exclude (may be repeated)
* `<docroot>`: document root directory to scan (defaults to the current working directory; may be repeated)

AWS credentials are sourced from the environment or AWS SDK configuration files. There are no command-line options for passing access key IDs, secrets, or session tokens.

See also
--------

Sitesearch is part of the [Mergician](https://github.com/rcrowley/mergician) suite of tools that manipulate HTML documents:

* [Deadlinks](https://github.com/rcrowley/deadlinks): Scan a document root directory for dead links
* [Electrostatic](https://github.com/rcrowley/electrostatic): Mergician-powered, pure-HTML CMS
* [Feed](https://github.com/rcrowley/feed): Scan a document root directory to construct an Atom feed
* [Frag](https://github.com/rcrowley/frag): Extract fragments of HTML documents
