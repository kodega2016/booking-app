package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type MyHandler struct{}

func (myH *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
