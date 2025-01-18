package routing

import "net/http"

type handler struct {
    handler func(w http.ResponseWriter, r *http.Request)
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.handler(w,r)
}

type Route struct {
    Path string
    Method string
    Handler http.Handler
}

func createRoute(method string, path string, handleFunc func(http.ResponseWriter, *http.Request), middlewares []Middleware) *Route {
    var handler http.Handler = handler{
        handler: handleFunc,
    }
    if len(middlewares)!=0{
        for _, middleware := range middlewares {
            handler = middleware(handler)
        }
    }

    route := Route{
        Path: path,
        Method: method,
        Handler: handler,
    }

    return &route
}

type Router struct {
    Name string
    Routes []*Route
}

func (r *Router) Get(path string, handler func(http.ResponseWriter, *http.Request), middlewares []Middleware) {
    route := createRoute("GET", path, handler, middlewares)
    r.Routes = append(r.Routes, route)
}

func (r *Router) Post(path string, handler func(http.ResponseWriter, *http.Request), middlewares []Middleware) {
    route := createRoute("POST", path, handler, middlewares)
    r.Routes = append(r.Routes, route)
}

func (r *Router) Put(path string, handler func(http.ResponseWriter, *http.Request), middlewares []Middleware) {
    route := createRoute("PUT", path, handler, middlewares)
    r.Routes = append(r.Routes, route)
}

func (r *Router) Patch(path string, handler func(http.ResponseWriter, *http.Request), middlewares []Middleware) {
    route := createRoute("PUT", path, handler, middlewares)
    r.Routes = append(r.Routes, route)
}

func (r *Router) Delete(path string, handler func(http.ResponseWriter, *http.Request), middlewares []Middleware) {
    route := createRoute("DELETE", path, handler, middlewares)
    r.Routes = append(r.Routes, route)
}
