package routing

import (
	"net/http"
	"strconv"
)

type Server struct {
    http.Server
    Protocol string
    FullAddr string
    Routers []*Router
}

type Middleware func(http.Handler) http.Handler

type Options struct {
    Host string
    Port int
    Protocol string
    Middlewares []Middleware
}


func (server *Server) createHandler(middlewares []Middleware ) (http.Handler) {
    mux := http.NewServeMux()

    for _, router := range server.Routers {
        base:=""
        
        if router.Name!=""{
            base = "/" + router.Name
        }

        for _, route := range router.Routes{
            mux.Handle(route.Method + " " + base + route.Path, route.Handler)
        }
    }

    var handler http.Handler = mux

    for _,middleware := range middlewares {
        handler = middleware(handler)
    }

    return handler
}

func (server *Server) AddRouter (router *Router){
    server.Routers = append(server.Routers, router)
}

func (server *Server) InitServer (options *Options) {
    handler := server.createHandler(options.Middlewares)

    protocol := "http"

    if options.Host!="" {
        protocol = options.Protocol
    }

    host := "localhost"

    if options.Host!="" {
        host = options.Host
    }

    port := 8080

    if options.Port!=0 {
        port = options.Port
    }

    address := host + ":" + strconv.Itoa(port)

    server.Server = http.Server {
        Addr: address,
        Handler: handler,
    }
    
    server.Protocol = protocol
    server.FullAddr = server.Protocol + "://" + server.Addr
}

