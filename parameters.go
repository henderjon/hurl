package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	help           bool
	stderr, stdout *log.Logger
)

type multiParams map[string][]string

type getOptParameters struct {
	optFormURLEncode bool
	optSilence       bool
	optQueryString   bool
	optPostForm      bool
	optOutFile       bool
	optReadStdin     bool
	optSummary       bool
	optPost          bool
	optHTTPAction    string
	optURI           string
	optBasic         string
	optToken         string
	optBearer        string
	optType          string
	optBinData       string
	optHeaders       multiParams
	optData          multiParams
}

const doc = `
%s is a tool for looking at the request being made

version:  %s
compiled: %s
built:    %s

Usage: %s -u <URL> [option [option]...]

Options:
`

// GetParams parses CLI args into values used by the application
func getParams(buildVersion, buildTimestamp, compiledBy string) *getOptParameters {
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			doc,
			os.Args[0],
			buildVersion,
			compiledBy,
			buildTimestamp,
			os.Args[0],
		)
		flag.PrintDefaults()
	}

	params := &getOptParameters{
		optHeaders: multiParams{},
		optData:    multiParams{},
	}
	// redis
	flag.BoolVar(&help, "help", false, "display these program options")
	flag.BoolVar(&params.optPost, "post", false, "sugar for setting the HTTP action to POST")
	flag.BoolVar(&params.optFormURLEncode, "form", false, "sugar for adding 'Content-Type: application/x-www-form-urlencoded'")
	flag.BoolVar(&params.optPostForm, "pf", false, "sugar for -post -form")
	flag.BoolVar(&params.optSilence, "s", false, "shutup")
	flag.BoolVar(&params.optQueryString, "query", false, "append -data-urlencoded's to the target URL as a query string")
	flag.BoolVar(&params.optOutFile, "save", false, "write the output to a similarly named local file; to specify a different filename, simply redirect stdout")
	flag.BoolVar(&params.optReadStdin, "stdin", false, "read the request body from stdin; request will ingore 'param' and 'body'")
	flag.BoolVar(&params.optSummary, "summary", false, "after the request is finished, print a brief summary")
	flag.Var(&params.optHeaders, "header", "`param=value` headers for the request")
	flag.Var(&params.optData, "param", "`param=value` data for the request")
	flag.StringVar(&params.optBinData, "body", "", "data as a string for the body of the request")
	flag.StringVar(&params.optHTTPAction, "X", "GET", "specify the HTTP `action` (e.g. GET, POST, etc)")

	flag.StringVar(&params.optURI, "url", "", "the destination URI")
	flag.StringVar(&params.optBasic, "basic", "", "sugar for adding the 'Authorization: Basic $val' header")
	flag.StringVar(&params.optToken, "token", "", "sugar for adding the 'Authorization: Token $val' header")
	flag.StringVar(&params.optBearer, "bearer", "", "sugar for adding the 'Authorization: Bearer $val' header")
	flag.StringVar(&params.optType, "type", "", "sugar for adding the 'Content-Type: $val' header")

	flag.Parse()

	if help || len(params.optURI) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	stdout = log.New(os.Stdout, "", 0)
	stderr = log.New(os.Stderr, "", 0)
	if params.optSilence {
		stderr = log.New(ioutil.Discard, "", 0)
	}

	return params
}
