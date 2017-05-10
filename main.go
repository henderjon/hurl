package main

import (
	"fmt"
	"log"
	"net/url"
	"net/http"
	"io"
	"os"
	"strings"
)

const (
	// escape = "\x1b"
	prefix_in = ""
	prefix_out = ""
	protocol = "\x1b[91m%s %s %s\x1b[0m\n"
	header = "%s\x1b[90m%s:\x1b[0m \x1b[94m%s\x1b[0m\n"
)


func main() {
	client := http.Client{}

	target := "httpbin.org/get"

	remote, err := url.Parse(target);
	if err != nil {
		log.Fatal(err)
	}

	if remote.Scheme == "" {
		remote.Scheme = "http"
		remote, _ = url.Parse(remote.String())
	}

	req, err := http.NewRequest("GET", remote.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Vnd.hurl.version 4")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", remote.Host)

	fmt.Printf(protocol, req.Method, req.URL.Path, req.Proto)
	for k, v := range req.Header {
		fmt.Printf(header, prefix_out, k, strings.Join(v, ", "))
	}

	fmt.Print("\n")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Printf(protocol, resp.Proto, resp.Status, "")

	for k, v := range resp.Header {
		fmt.Printf(header, prefix_in, k, strings.Join(v, ", "))
	}

	fmt.Print("\n")

	io.Copy(os.Stdout, resp.Body)
}
