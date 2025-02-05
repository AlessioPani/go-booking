package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	mux := routes()

	switch v := mux.(type) {
	case *chi.Mux:
	// do nothing
	default:
		t.Errorf("Type is not *chi.Mux, but %T", v)
	}

}
