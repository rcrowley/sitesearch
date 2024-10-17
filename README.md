Sitesearch
==========

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
sitesearch -l <layout> [-n <name>] [-r <region>] <input>[...]
```

* `-l <layout>`: site layout HTML document for search result pages
* `-n <name>`: name of the the Lambda function (default "sitesearch")
* `-r <region>`: AWS region to host the Lambda function (default to AWS\_DEFAULT\_REGION in the environment)
* `<input>[...]`: pathname, relative to your site's root, of one or more HTML files, given as command-line arguments or on standard input

See also
--------

Sitesearch is part of the [Mergician](https://github.com/rcrowley/mergician) suite of tools that manipulate HTML documents:

* [Deadlinks](https://github.com/rcrowley/deadlinks): Scan a document root directory for dead links
* [Electrostatic](https://github.com/rcrowley/electrostatic): Mergician-powered, pure-HTML CMS
* [Frag](https://github.com/rcrowley/frag): Extract fragments of HTML documents
