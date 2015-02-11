package router

import (
	"net/http"
	"regexp"
)

// Route struct
type Route struct {
	name    string         // Name of route
	pattern *regexp.Regexp // Pattern of route
	method  string         // Method of route: GET/POST/PUT/DELETE
	handler http.Handler   // HTTP Handler of route
}

func (this *Route) Name() string {
	return this.name
}

func (this *Route) Method() string {
	return this.method
}

func (this *Route) Pattern() *regexp.Regexp {
	return this.pattern
}

func (this *Route) Match(r *http.Request) bool {
	return (this.method == "" || this.method == r.Method) &&
		(this.pattern == nil || this.pattern.MatchString(r.URL.Path))
}

// Router struct
type Router struct {
	routes []*Route // All routes
}

func (this *Router) Get(pattern string, handler http.Handler) *Router {
	return this.AddRoute("GET", pattern, handler)
}

func (this *Router) Post(pattern string, handler http.Handler) *Router {
	return this.AddRoute("POST", pattern, handler)
}

func (this *Router) Put(pattern string, handler http.Handler) *Router {
	return this.AddRoute("PUT", pattern, handler)
}

func (this *Router) Delete(pattern string, handler http.Handler) *Router {
	return this.AddRoute("DELETE", pattern, handler)
}

func (this *Router) Use(handler http.Handler) *Router {
	return this.AddRoute("", "", handler)
}

func (this *Router) Handle(pattern string, handler http.Handler) *Router {
	return this.AddRoute("", pattern, handler)
}

func (this *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) *Router {
	return this.AddRoute("", pattern, http.HandlerFunc(handler))
}

func (this *Router) AddRoute(method, pattern string, handler http.Handler) *Router {
	route := new(Route)
	route.pattern = regexp.MustCompile(pattern)
	route.method = method
	route.handler = handler
	this.routes = append(this.routes, route)
	return this
}

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range this.routes {
		if route.Match(r) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
}

// func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	for _, route := range this.routes {
// 		if route.Match(r) {
// 			println(r.URL.Path)
// 			eachHandler(w, r, route.handlers)
// 			return
// 		}
// 	}
// 	// no pattern matched; send 404 response
// 	eachHandler(w, r, this.notFound)
// }
//
// func eachHandler(w http.ResponseWriter, r *http.Request, handlers []Handler) {
// 	l := len(handlers)
// 	if l < 1 {
// 		return
// 	} else if l == 1 {
// 		handlers[0](w, r, func() {})
// 		return
// 	}
// 	handlers[0](w, r, func() {
// 		eachHandler(w, r, handlers[1:])
// 	})
// }
