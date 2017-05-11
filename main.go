package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

const (
	// escape = "\x1b"
	nl        = "\n"
	prefixIn  = ""
	prefixOut = ""
	protocol  = "\x1b[91m%s %s %s\x1b[0m\n"
	header    = "%s\x1b[90m%s:\x1b[0m \x1b[94m%s\x1b[0m\n"
	userAgent = "henderjon/hurl v1.0"
)

var (
	optFormURLEncode, optSilence,
	optQueryString, optPostForm,
	optOutFile, help bool
	optHTTPAction, optTargetURL string
	stderr, stdout              *log.Logger
	optHeaders                  = multiParams{}
	optData                     = multiParams{}
)

func init() {

	flag.BoolVar(&help, "help", false, "display these program options")
	flag.BoolVar(&optFormURLEncode, "f", false, "sugar for adding 'Content-Type: application/x-www-form-urlencoded'")
	flag.BoolVar(&optPostForm, "pf", false, "form-urlencode the POST; sugar for '-X POST -f'")
	flag.BoolVar(&optSilence, "s", false, "shutup")
	flag.BoolVar(&optQueryString, "q", false, "append -d's to the target URL as a query string")
	flag.BoolVar(&optOutFile, "save", false, "write the output to a similarly named local file; to specify a different filename, simply redirect stdout")
	flag.Var(&optHeaders, "h", "`param=value` headers for the request")
	flag.Var(&optData, "d", "`param=value` data for the request")
	flag.StringVar(&optHTTPAction, "X", "GET", "specify the HTTP `action` (e.g. GET, POST, etc)")

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		log.Fatal("The URL should be the last arg")
	}

	optTargetURL = args[len(args)-1]

	stdout = log.New(os.Stdout, "", 0)
	stderr = log.New(os.Stderr, "", 0)
	if optSilence {
		stderr = log.New(ioutil.Discard, "", 0)
	}

}

func main() {

	if help {
		fmt.Println("\n\x1b[91mhurl is a utility for making HTTP requests. \x1b[0m\n")
		flag.PrintDefaults()
		return
	}

	client := http.Client{}

	remote, err := url.Parse(optTargetURL)
	if err != nil {
		log.Fatal("Unable to parse", optTargetURL)
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

	if optPostForm {
		optHTTPAction = http.MethodPost
		optFormURLEncode = true
	}

	var body io.Reader
	if optQueryString {
		remote.RawQuery = data.Encode() // force a query string with -q
	} else {
		body = bytes.NewBufferString(data.Encode()) // send the data as the body
	}

	req, err := http.NewRequest(optHTTPAction, remote.String(), body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", remote.Host)
	if optFormURLEncode {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	for name, slice := range optHeaders {
		for _, v := range slice {
			req.Header.Add(name, v)
		}
	}

	stderr.Printf(protocol, req.Method, req.URL.RequestURI(), req.Proto)
	for k, v := range req.Header {
		for _, v := range v {
			stderr.Printf(header, prefixOut, k, v)
		}
	}

	stderr.Print("\n")
	stderr.Printf("%v\n", body)
	stderr.Print("\n")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	stderr.Printf(protocol, resp.Proto, resp.Status, "")

	printHeaders(resp.Header)

	stderr.Print("\n")

	local := os.Stdout
	if optOutFile {
		fname := path.Base(remote.Path)
		local, err = os.Open(fname)
		if os.IsNotExist(err) {
			local, _ = os.Create(fname)
		} else {
			log.Fatal(err)
		}
	}

	io.Copy(local, resp.Body)
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
	v := strings.Split(value, "=")
	if len(v) == 2 {
		d[v[0]] = append(d[v[0]], v[1])
	}
	return nil
}
