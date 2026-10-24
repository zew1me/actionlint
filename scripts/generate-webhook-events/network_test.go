//go:build network

// Tests in this file need outbound network access. Run them with `go test -tags network`.

package main

import "testing"

func TestFetchURL(t *testing.T) {
	b, err := fetch("https://github.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatal("Fetched source is empty")
	}
}
