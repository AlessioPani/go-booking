package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// Test setup

	// After testing setup, start the actual tests
	os.Exit(m.Run())
}

type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
