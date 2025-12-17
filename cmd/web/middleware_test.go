package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myHandler MyHandler
	h := NoSurf(&myHandler)

	switch v := h.(type) {
	case http.Handler:
	default:
		t.Errorf("type is not http.Handler,we got %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myHandler MyHandler
	h := SessionLoad(&myHandler)

	switch v := h.(type) {
	case http.Handler:
	default:
		t.Errorf("type is not http.Handler,we got %T", v)
	}
}
