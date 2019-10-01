[![GoDoc](https://godoc.org/github.com/henderjon/hurl?status.svg)](https://godoc.org/github.com/henderjon/hurl)
[![License: BSD-3](https://img.shields.io/badge/license-BSD--3-blue.svg)](https://img.shields.io/badge/license-BSD--3-blue.svg)
![tag](https://img.shields.io/github/tag/henderjon/hurl.svg)
![release](https://img.shields.io/github/release/henderjon/hurl.svg)

[![Go Report Card](https://goreportcard.com/badge/github.com/henderjon/hurl)](https://goreportcard.com/report/github.com/henderjon/hurl)
[![Build Status](https://travis-ci.org/henderjon/hurl.svg?branch=dev)](https://travis-ci.org/henderjon/hurl)
[![Maintainability](https://api.codeclimate.com/v1/badges/df165f1d091666a37b09/maintainability)](https://codeclimate.com/github/henderjon/hurl/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/df165f1d091666a37b09/test_coverage)](https://codeclimate.com/github/henderjon/hurl/test_coverage)

## hurl is a tool for looking at the request being made

Inspired by [bat](https://github.com/astaxie/bat) and [kurly](https://github.com/davidjpeacock/kurly), `hurl` is simple HTTP requester.

```
Usage: ./hurl -u <URL> [option [option]...]

Options:
  -X action
    	specify the HTTP action (e.g. GET, POST, etc) (default "GET")
  -basic string
    	sugar for adding the 'Authorization: Basic $val' header
  -bearer string
    	sugar for adding the 'Authorization: Bearer $val' header
  -body string
    	data as a string for the body of the request
  -form
    	sugar for adding 'Content-Type: application/x-www-form-urlencoded'
  -header param=value
    	param=value headers for the request
  -help
    	display these program options
  -param param=value
    	param=value data for the request
  -pf
    	sugar for -post -form
  -post
    	sugar for setting the HTTP action to POST
  -query
    	append -data-urlencoded's to the target URL as a query string
  -s	shutup
  -save
    	write the output to a similarly named local file; to specify a different filename, simply redirect stdout
  -stdin
    	read the request body from stdin; request will ingore 'param' and 'body'
  -summary
    	after the request is finished, print a brief summary
  -token string
    	sugar for adding the 'Authorization: Token $val' header
  -type string
    	sugar for adding the 'Content-Type: $val' header
  -url string
    	the destination URI
```

### why

Whenever I build HTTP APis, there seems to be a number of utilities that do more than I want or do what I want in ways that are more complex than I would prefer. This is my attempt at a utility that is small, simple, and does things the way I would do them. Where I tend to spend my keystrokes, I've added sugar.

### todo

  - progress bars
  - multipart/form-data (binary data)

