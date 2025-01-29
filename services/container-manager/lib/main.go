//package main
//
//import (
//	"fmt"
//	//"github.com/DraouiBilal/code-runner-backend-lib/api"
//	"github.com/DraouiBilal/code-runner-backend-lib/routing"
//	"log"
//	"net/http"
//)
//
//func main() {
//	server := routing.Server{}
//	dockerRouter := &routing.Router{
//		Name: "docker",
//	}
//
//	dockerRouter.Get("/test/{id}", func(w http.ResponseWriter, req *http.Request) {
//		routing.Utils.WriteJSON(w, struct{ Test string }{Test: req.PathValue("id")})
//	    w.WriteHeader(http.StatusCreated)
//	}, []routing.Middleware{})
//
//	server.AddRouter(dockerRouter)
//
//	server.InitServer(&routing.Options{
//		Middlewares: []routing.Middleware{routing.Middlewares.Logging},
//	})
//
//	log.Printf("Server starting on %s", server.FullAddr)
//
//	err := server.ListenAndServe()
//
//	if err != nil {
//		fmt.Println(err)
//	}
//}
