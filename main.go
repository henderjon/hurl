package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	// escape = "\x1b"
	nl        = "\n"
	prefixIn  = ""
	prefixOut = ""
	protocol  = "\x1b[91m%s %s %s\x1b[0m\n"
	header    = "%s\x1b[90m%s:\x1b[0m \x1b[94m%s\x1b[0m\n"
	summary   = "\x1b[90m%s:\x1b[0m \x1b[94m%s\x1b[0m\n"
	userAgent = "hurl/v0.1.3-alpha"
)

var (
	optFormURLEncode, optSilence,
	optQueryString, optPostForm,
	optOutFile, optReadStdin,
	optSummary, optPost, help bool
	optHTTPAction, optURI, optBasic,
	optToken, optBearer, optType,
	optBinData string
	stderr, stdout *log.Logger
	optHeaders     = multiParams{}
	optData        = multiParams{}
)

func init() {

	flag.BoolVar(&help, "help", false, "display these program options")
	flag.BoolVar(&optPost, "post", false, "set the HTTP action to POST; this is sugar")
	flag.BoolVar(&optFormURLEncode, "f", false, "sugar for adding 'Content-Type: application/x-www-form-urlencoded'")
	flag.BoolVar(&optPostForm, "pf", false, "form-urlencode the POST; sugar for '-X POST -f'")
	flag.BoolVar(&optSilence, "s", false, "shutup")
	flag.BoolVar(&optQueryString, "q", false, "append -d's to the target URL as a query string")
	flag.BoolVar(&optOutFile, "save", false, "write the output to a similarly named local file; to specify a different filename, simply redirect stdout")
	flag.BoolVar(&optReadStdin, "stdin", false, "read the request body from stdin; request will ingore all -d's")
	flag.BoolVar(&optSummary, "summary", false, "after the request is finished, print a brief summary")
	flag.Var(&optHeaders, "h", "`param=value` headers for the request")
	flag.Var(&optData, "d", "`param=value` data for the request")
	flag.StringVar(&optBinData, "bin", "", "data as a string for the body of the request")
	flag.StringVar(&optHTTPAction, "X", "GET", "specify the HTTP `action` (e.g. GET, POST, etc)")

	flag.StringVar(&optURI, "u", "", "the destination URI; if not provided the URI is assumed to be the last arg")
	flag.StringVar(&optBasic, "basic", "", "sugar for adding the 'Authorization: Basic $val' header")
	flag.StringVar(&optToken, "token", "", "sugar for adding the 'Authorization: Token $val' header")
	flag.StringVar(&optBearer, "bearer", "", "sugar for adding the 'Authorization: Bearer $val' header")
	flag.StringVar(&optType, "type", "", "sugar for adding the 'Content-Type: $val' header")

	flag.Parse()

	if help {
		fmt.Println("\n\x1b[91mhurl is a utility for making HTTP requests. \x1b[0m\n")
		flag.PrintDefaults()
		os.Exit(0)
	}

	args := flag.Args()

	if len(optURI) == 0 {
		if len(args) == 0 {
			log.Fatal("The URL should be the last arg or should be specified by -u; use -help for more info")
		}

		optURI = args[len(args)-1]
	}

	stdout = log.New(os.Stdout, "", 0)
	stderr = log.New(os.Stderr, "", 0)
	if optSilence {
		stderr = log.New(ioutil.Discard, "", 0)
	}

}

func main() {

	client := http.Client{}

	remote, err := url.Parse(optURI)
	if err != nil {
		log.Fatal("Unable to parse", optURI)
	}

	if remote.Scheme == "" {
		remote.Scheme = "http"
		remote, _ = url.Parse(remote.String()) // reparse to populate remote.Host
	}

	data := url.Values{}
	for name, slice := range optData {
		for _, v := range slice {
			data.Add(name, v)
		}
	}

	switch {
	case optPostForm:
		optFormURLEncode = true
		fallthrough
	case optPost:
		optHTTPAction = http.MethodPost
	}

	var body bytes.Buffer // io.ReadWriter
	switch {
	case optQueryString:
		remote.RawQuery = data.Encode() // force a query string with -q
	case optReadStdin:
		if optReadStdin {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				body.Write(scanner.Bytes())
			}
		}
	case len(optBinData) > 0:
		body.WriteString(optBinData)
	case body.Len() == 0:
		fallthrough
	default:
		body.WriteString(data.Encode()) // send the query data as the body
	}

	req, err := http.NewRequest(optHTTPAction, remote.String(), &body)
	if err != nil {
		log.Fatal(err)
	}

	// sugar for basic auth
	if len(optBasic) > 0 {
		if strings.Contains(optBasic, ":") {
			optBasic = base64.StdEncoding.EncodeToString([]byte(optBasic))
		}
		req.Header.Set("Authorization", "Basic "+optBasic)
	}

	// sugar for token auth
	if len(optToken) > 0 {
		req.Header.Set("Authorization", "Token "+optToken)
	}

	// sugar for bearer auth
	if len(optBearer) > 0 {
		req.Header.Set("Authorization", "Bearer "+optBearer)
	}

	// sugar for Content-Type
	if len(optType) > 0 {
		req.Header.Set("Content-Type", optType)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", remote.Host)
	req.Header.Set("Content-Length", strconv.Itoa(body.Len()))
	if optFormURLEncode {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	}

	for name, slice := range optHeaders {
		for _, v := range slice {
			req.Header.Add(name, v)
		}
	}

	stderr.Printf(protocol, req.Method, req.URL.RequestURI(), req.Proto)
	printHeaders(req.Header)

	stderr.Print(nl)
	stderr.Printf("%v%s", &body, nl)
	stderr.Print(nl)

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	stderr.Printf(protocol, resp.Proto, resp.Status, "")

	printHeaders(resp.Header)

	stderr.Print(nl)

	local := os.Stdout
	if optOutFile {
		// log.Fatal(resp.Header.Get("Content-Disposition"))
		fname := path.Base(remote.Path)
		local, err = os.OpenFile(fname, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	n, err := io.Copy(local, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	stderr.Println(nl)

	if optSummary {
		// stderr.Println("\x1b[91m//------------------------------------------------------------------------//\x1b[0m")
		stderr.Println("\x1b[91m---------------\x1b[0m")
		stderr.Printf(summary, "Content-Length", strconv.FormatInt(resp.ContentLength, 10))
		stderr.Printf(summary, "Bytes Received", strconv.FormatInt(n, 10))
		stderr.Printf(summary, "Request Duration", time.Since(start).String())
	}
}

func printHeaders(headers map[string][]string) {
	for k, v := range headers {
		for _, v := range v {
			stderr.Printf(header, prefixOut, k, v)
		}
	}
}

type multiParams map[string][]string

func (d multiParams) String() string {
	return fmt.Sprintf("%d", d)
}

func (d multiParams) Set(value string) error {
	v := strings.SplitN(value, "=", 2)
	if len(v) == 2 {
		d[v[0]] = append(d[v[0]], v[1])
	}
	return nil
}
