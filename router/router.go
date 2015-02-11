package router

import (
	"net/http"
	"regexp"
)

// Route struct
type Route struct {
	name    string         // Name of route
	regexp  *regexp.Regexp // Regexp for route
	pattern string         // Pattern of route
	method  string         // Method of route: GET/POST/PUT/DELETE
	handler http.Handler   // HTTP Handler of route
}

func (this *Route) Name() string {
	return this.name
}

func (this *Route) Method() string {
	return this.method
}

func (this *Route) Pattern() string {
	return this.pattern
}

func (this *Route) MatchMethod(method string) bool {
	return this.method == "" || this.method == method || (method == "HEAD" && this.method == "GET")
}

func (this *Route) Match(r *http.Request) bool {
	return this.MatchMethod(r.Method) && this.regexp.MatchString(r.URL.Path)
}

// Router struct
type Router struct {
	routes   []*Route     // All routes
	notFound http.Handler // not found handler
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

func (this *Router) NotFound(handler http.Handler) *Router {
	this.notFound = handler
	return this
}

func (this *Router) AddRoute(method, pattern string, handler http.Handler) *Router {
	route := new(Route)
	route.pattern = pattern
	route.method = method
	route.handler = handler
	route.regexp = regexp.MustCompile(pattern)
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
	if this.notFound != nil {
		this.notFound.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}
