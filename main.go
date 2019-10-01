package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
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
)

func main() {

	params := getParams(getBuildVersion(), getBuildTimestamp(), getCompiledBy())

	client := http.Client{}

	remote, err := url.Parse(params.optURI)
	if err != nil {
		log.Fatal("Unable to parse", params.optURI)
	}

	if remote.Scheme == "" {
		remote.Scheme = "http"
		remote, _ = url.Parse(remote.String()) // reparse to populate remote.Host
	}

	data := url.Values{}
	parseMultiData(params.optData, data)

	switch {
	case params.optPostForm:
		params.optFormURLEncode = true
		fallthrough
	case params.optPost:
		params.optHTTPAction = http.MethodPost
	}

	var body bytes.Buffer // io.ReadWriter
	switch {
	case params.optQueryString:
		remote.RawQuery = data.Encode() // force a query string with -q
	case params.optReadStdin:
		if params.optReadStdin {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				body.Write(scanner.Bytes())
			}
		}
	case len(params.optBinData) > 0:
		body.WriteString(params.optBinData)
	case body.Len() == 0:
		fallthrough
	default:
		body.WriteString(data.Encode()) // send the query data as the body
	}

	req, err := http.NewRequest(params.optHTTPAction, remote.String(), &body)
	if err != nil {
		log.Fatal(err)
	}

	// sugar for basic auth
	if len(params.optBasic) > 0 {
		if strings.Contains(params.optBasic, ":") {
			params.optBasic = base64.StdEncoding.EncodeToString([]byte(params.optBasic))
		}
		req.Header.Set("Authorization", "Basic "+params.optBasic)
	}

	// sugar for token auth
	if len(params.optToken) > 0 {
		req.Header.Set("Authorization", "Token "+params.optToken)
	}

	// sugar for bearer auth
	if len(params.optBearer) > 0 {
		req.Header.Set("Authorization", "Bearer "+params.optBearer)
	}

	// sugar for Content-Type
	if len(params.optType) > 0 {
		req.Header.Set("Content-Type", params.optType)
	}

	req.Header.Set("User-Agent", "hurl/"+buildVersion)
	req.Header.Set("Accept", "*/*")
	// req.Header.Set("cache-control", "no-cache")
	req.Header.Set("Host", remote.Host)

	if body.Len() > 0 {
		req.Header.Set("Content-Length", strconv.Itoa(body.Len()))
	}

	if params.optFormURLEncode {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	}

	parseMultiData(params.optHeaders, req.Header)

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
	if params.optOutFile {
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

	if params.optSummary {
		// stderr.Println("\x1b[91m//------------------------------------------------------------------------//\x1b[0m")
		stderr.Println("\x1b[91m---------------\x1b[0m")
		stderr.Printf(summary, "Content-Length", strconv.FormatInt(resp.ContentLength, 10))
		stderr.Printf(summary, "Bytes Received", strconv.FormatInt(n, 10))
		stderr.Printf(summary, "Request Duration", time.Since(start).String())
	}
}

type adder interface {
	Add(string, string)
}

func parseMultiData(src multiParams, dest adder) {
	for name, slice := range src {
		for _, v := range slice {
			dest.Add(name, v)
		}
	}
}

func printHeaders(headers map[string][]string) {
	for k, v := range headers {
		for _, v := range v {
			stderr.Printf(header, prefixOut, k, v)
		}
	}
}

func (d multiParams) String() string {
	return fmt.Sprintf("%d", len(d))
}

func (d multiParams) Set(value string) error {
	v := strings.SplitN(value, "=", 2)
	if len(v) == 2 {
		d[v[0]] = append(d[v[0]], v[1])
	}
	return nil
}
