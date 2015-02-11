package handlers

import (
	"net/http"

	"github.com/jdkanani/smalldocs/context"
)

/**
 * Handle function
 */
type HandleFunc func(*context.Context, http.ResponseWriter, *http.Request) (int, error)

/**
 *  Handle interace
 */
type Handler interface {
	ServeHTTP(*context.Context, http.ResponseWriter, http.Request) (int, error)
}
