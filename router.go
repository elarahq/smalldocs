package main

import (
	"net/http"
	"regexp"
)

// Http handler
type Next func()
type Handler func(w http.ResponseWriter, r *http.Request, next Next)

// Route struct
type Route struct {
	name     string         // Name of route
	pattern  *regexp.Regexp // Pattern of route
	method   string         // Method of route: GET/POST/PUT/DELETE
	handlers []Handler      // HTTP Handlers of route
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
	routes       []*Route  // All routes
	notFound     []Handler // Handlers for "NOT FOUND"
	errorOccured []Handler // Handlers for "ERROR"
}

func (this *Router) Get(pattern string, handlers ...Handler) *Router {
	return this.AddRoute("GET", pattern, handlers...)
}

func (this *Router) Post(pattern string, handlers ...Handler) *Router {
	return this.AddRoute("POST", pattern, handlers...)
}

func (this *Router) Put(pattern string, handlers ...Handler) *Router {
	return this.AddRoute("PUT", pattern, handlers...)
}

func (this *Router) Delete(pattern string, handlers ...Handler) *Router {
	return this.AddRoute("DELETE", pattern, handlers...)
}

func (this *Router) Use(handlers ...Handler) *Router {
	return this.AddRoute("", "", handlers...)
}

func (this *Router) Handler(pattern string, handlers ...Handler) *Router {
	return this.AddRoute("", pattern, handlers...)
}

func (this *Router) HandlerFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) *Router {
	wrapper := func(w http.ResponseWriter, r *http.Request, next Next) {
		handler(w, r)
	}
	return this.AddRoute("", pattern, wrapper)
}

func (this *Router) AddRoute(method, pattern string, handlers ...Handler) *Router {
	route := new(Route)

	route.pattern = regexp.MustCompile(pattern)
	route.method = method

	route.handlers = make([]Handler, 0)
	route.handlers = append(route.handlers, handlers...)

	this.routes = append(this.routes, route)
	return this
}

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range this.routes {
		if route.Match(r) {
			println(r.URL.Path)
			eachHandler(w, r, route.handlers)
			return
		}
	}
	// no pattern matched; send 404 response
	eachHandler(w, r, this.notFound)
}

func eachHandler(w http.ResponseWriter, r *http.Request, handlers []Handler) {
	l := len(handlers)
	if l < 1 {
		return
	} else if l == 1 {
		handlers[0](w, r, func() {})
		return
	}
	handlers[0](w, r, func() {
		eachHandler(w, r, handlers[1:])
	})
}
