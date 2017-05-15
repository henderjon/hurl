package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"
)

const fmtFieldname = "Content-Disposition: form-data; name=\"%s\"\n"
const fmtFilename = "Content-Disposition: form-data; name=\"%s\"; filename=\"%s\"\n"
const fmtLength = "Content-Length: %s\n\n"
const fmtMimetype = "Content-Type: %s\n"

func parseBinData(body *bytes.Buffer, binData multiParams) {
	boundaryStr := "--" + boundary + eol
	for name, slice := range binData {
		for _, v := range slice {
			if f, err := ioutil.ReadFile(v); err == nil {
				body.WriteString(boundaryStr)
				body.WriteString(fmt.Sprintf(fmtFilename, name, path.Base(v)))
				body.WriteString(fmt.Sprintf(fmtMimetype, http.DetectContentType(f)))
				body.WriteString(fmt.Sprintf(fmtLength, strconv.Itoa(len(f))))
				body.Write(f)
			} else {
				log.Fatal(err)
			}
		}
	}
}

func parseData(body *bytes.Buffer, binData multiParams) {
	boundaryStr := "--" + boundary + eol
	for name, slice := range binData {
		for _, v := range slice {
			body.WriteString(boundaryStr)
			body.WriteString(fmt.Sprintf(fmtFieldname, name))
			body.WriteString(fmt.Sprintf(fmtLength, strconv.Itoa(len(v))))
			body.WriteString(v)
			body.WriteString(eol)
		}
	}
}
