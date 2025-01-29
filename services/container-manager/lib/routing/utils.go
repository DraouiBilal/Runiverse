package routing

import (
	"encoding/json"
	"log"
	"net/http"
)


type utils struct {
    WriteJSON func(http.ResponseWriter, interface{})
}

func WriteJSON(w http.ResponseWriter, body interface{}) {

    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

    payload, err := json.Marshal(body) 

    if err != nil {
        log.Fatal(err)
    }

    _, err = w.Write(payload)
    
    if err != nil {
        log.Fatal(err)
    }
}

var Utils *utils = &utils{
    WriteJSON: WriteJSON,
}
