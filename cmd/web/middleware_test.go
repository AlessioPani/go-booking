package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var mH myHandler
	h := NoSurf(&mH)

	switch v := h.(type) {
	case http.Handler:
	// do nothing
	default:
		t.Errorf("Type is not an http.Handler %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var mH myHandler
	h := SessionLoad(&mH)

	switch v := h.(type) {
	case http.Handler:
	// do nothing
	default:
		t.Errorf("Type is not an http.Handler %T", v)
	}

}
