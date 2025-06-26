package server

import (
	"fmt"
	"net/http"
	"strings"
)

type Handler func(w http.ResponseWriter, r *http.Request)

type Router struct {
	routes          map[string]map[string]Handler
	server          *Server
	notFoundHandler Handler
}

func newRouter(serv *Server) *Router {
	return &Router{
		routes:          make(map[string]map[string]Handler),
		server:          serv,
		notFoundHandler: defaultNotFoundHandler,
	}
}

func (router *Router) notFound(handler Handler) {
	router.notFoundHandler = handler
}

func (router *Router) addRoute(method, path string, handler Handler) {
	if router.routes[method] == nil {
		router.routes[method] = make(map[string]Handler)
	}
	router.routes[method][path] = handler
}

func (router *Router) GET(path string, handler Handler) {
	router.addRoute(http.MethodGet, path, handler)
}

func (router *Router) POST(path string, handler Handler) {
	router.addRoute(http.MethodPost, path, handler)
}

func (router *Router) PUT(path string, handler Handler) {
	router.addRoute(http.MethodPut, path, handler)
}

func (router *Router) DELETE(path string, handler Handler) {
	router.addRoute(http.MethodDelete, path, handler)
}

func (router *Router) serveHTTP(w http.ResponseWriter, r *http.Request) {
	handler, err := router.findHandler(r.Method, r.URL.Path)
	if err != nil {
		router.notFoundHandler(w, r)
		return
	}
	if router.server != nil {
		handler = router.server.applyMiddleware(handler)
	}
	handler(w, r)
}

func (router *Router) findHandler(method, path string) (Handler, error) {
	if methodRoutes, ok := router.routes[method]; ok {
		if handler, ok := methodRoutes[path]; ok {
			return handler, nil
		}

		for routePath, handler := range methodRoutes {
			if isWildcardMatch(routePath, path) {
				return handler, nil
			}
		}
	}
	return nil, fmt.Errorf("no handler found for %s %s", method, path)
}

func isWildcardMatch(routePath, requestPath string) bool {
	routeParts := strings.Split(routePath, "/")
	requestParts := strings.Split(requestPath, "/")

	if len(routeParts) != len(requestParts) {
		return false
	}

	for i, part := range routeParts {
		if part == "*" {
			continue
		}
		if part != requestParts[i] {
			return false
		}
	}
	return true
}

func defaultNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 page not found", http.StatusNotFound)
}
