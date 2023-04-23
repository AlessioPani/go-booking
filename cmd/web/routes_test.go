package main

import (
	"fmt"
	"testing"

	"github.com/AlessioPani/go-booking/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
	// do nothing
	default:
		t.Errorf(fmt.Sprintf("Type is not *chi.Mux, but %T", v))
	}

}
