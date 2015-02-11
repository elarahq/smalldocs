package handlers

import (
	"net/http"

	"github.com/jdkanani/smalldocs/context"
)

type HandleFunc func(*context.Context, http.ResponseWriter, *http.Request) (int, error)

type Handler interface {
	ServeHTTP(*context.Context, http.ResponseWriter, http.Request) (int, error)
}
