package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jdkanani/smalldocs/context"
)

const (
	applicationJson = "application/json"
	textPlain       = "text/plain"
	textHTML        = "text/html"
)

//
// Error messages
//
var ErrorMessages = map[int]string{
	301: "Moved Permanently",
	304: "Not Modified",

	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	407: "Proxy Authentication Required",
	412: "Precondition Failed",

	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
}

//
// Error response
//
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (this ErrorResponse) Error() string {
	return fmt.Sprintf("%d %v", this.Code, this.Message)
}

//
// Error Handler
//
type ErrorHandler struct {
	Context *context.Context
	Handler HandleFunc
}

//
// Error handler type will now satisify http.Handler
// Errorhandler code
//
func (this *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ctx = this.Context
	if status, err := this.Handler(ctx, w, r); err != nil || status >= 400 {
		// set default status to 500
		if err != nil && status < 400 {
			status = http.StatusInternalServerError
		}

		// debug mode
		isDebug := ctx.Config.Bool("app.debug")
		// log errors
		if isDebug {
			log.Printf("HTTP %d: %v\n", status, r.URL.Path)
			if err != nil {
				fmt.Println(err)
			}
		}

		msg, ok := ErrorMessages[status]
		if !ok {
			status = http.StatusInternalServerError
			msg, _ = ErrorMessages[status]
		}

		em := ErrorResponse{
			Code:    status,
			Message: msg,
		}

		var e error
		// accept type
		accept := r.Header.Get("Accept")
		switch {
		case accept == applicationJson:
		case ctx.IsAjax(r):
			e = ctx.ErrorJSON(w, em, status)
		default:
			w.WriteHeader(status)
			e = ctx.RenderTemplate(w, "error", em.Error())
		}

		if e != nil {
			http.Error(w, "Error while rendering error page (Inception)!", http.StatusInternalServerError)
		}
	}
}
