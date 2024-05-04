package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

var helpErr = fmt.Errorf("usage: %s [PORT]", os.Args[0])

func run() error {
	if len(os.Args) != 2 {
		return helpErr
	}

	http.HandleFunc("/", handleProxy)
	return http.ListenAndServe(":"+os.Args[1], nil)
}

func handleProxy(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Handling request for: ", req.URL.EscapedPath())

	keyword := req.URL.Query().Get("keyword")

	path := req.URL.EscapedPath()
	path = strings.TrimLeft(path, "/")
	url := "https://" + path

	r, err := http.NewRequest(req.Method, url, req.Body)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "error proxying request: %s", err)
		fmt.Printf("error proxying request: %s", err)
		return
	}

	copyHeaders(req.Header, r.Header)
	client := http.Client{}

	res, err := client.Do(r)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "error proxying request: %s", err)
		fmt.Printf("error proxying request: %s", err)
		return
	}

	w.WriteHeader(res.StatusCode)

	copyHeaders(res.Header, w.Header())

	b := bufio.NewReader(res.Body)
	filterFeed(b, w, keyword)

	fmt.Println("Finished handling request")
}

func copyHeaders(src, dst http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func filterFeed(r byteReader, w io.Writer, keyword string) error {
	b := bytes.Buffer{}
	// Deliberately omit the closing '>' since some feeds have attributes
	// in this tag
	startDelim := []byte("<item")
	endDelim := []byte("</item>")

	for {
		err := search(r, w, startDelim)
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		err = search(r, &b, endDelim)
		if err != nil {
			return err
		}

		// There is a bit of memory copying going on here, but its
		// shouldn't really matter than much
		if strings.Contains(strings.ToLower(b.String()), keyword) {
			w.Write(startDelim)
			w.Write(b.Bytes())
			w.Write(endDelim)
		}
		b.Reset()
	}
}
