package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"reflect"
)

type Options struct {
	headers map[string]string
	Socket  string
}

func createRequest(url string, method string, body interface{}, options Options) *http.Request {
	payload, err := json.Marshal(body)

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "applciation/json")

	for _, k := range options.headers {
		req.Header.Set(k, options.headers[k])
	}

	if err != nil {
		log.Fatal(err)
	}

	return req
}

func makeRequest[T any](url string, method string, body interface{}, options Options) *T {
	var client *http.Client

	if options.Socket != "" {
		transport := &http.Transport{
			DialContext: func(context.Context, string, string) (net.Conn, error) {
				return net.Dial("unix", options.Socket)
			},
		}
		client = &http.Client{Transport: transport}
	} else {
		client = &http.Client{}
	}

	req := createRequest(url, method, body, options)

	res, err := client.Do(req)

	var result T

	switch any(result).(type) {
	case string:
		resText, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		strResult := string(resText)
		return any(&strResult).(*T)
	}

	if reflect.TypeOf((*T)(nil)).Elem().Kind() == reflect.Interface {
		return nil
	}

	if err != nil || res.StatusCode >= 400 {
        log.Println("Request failed or Status code >= 400")
		responseDump, err := httputil.DumpResponse(res, true)
		if err != nil {
			panic(err)
		}
		panic(string(responseDump))
	}

	defer res.Body.Close()

	responseJson := new(T)

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&responseJson)

	if err != nil {
        log.Println("JSON decoder failed")
		log.Fatal(err)
	}

	return responseJson
}

func Get[T any](url string, body interface{}, options Options) *T {
	return makeRequest[T](url, http.MethodGet, body, options)
}

func Post[T any](url string, body interface{}, options Options) *T {
	return makeRequest[T](url, http.MethodPost, body, options)
}

func Put[T any](url string, body interface{}, options Options) *T {
	return makeRequest[T](url, http.MethodPut, body, options)
}

func Patch[T any](url string, body interface{}, options Options) *T {
	return makeRequest[T](url, http.MethodPatch, body, options)
}

func Delete[T any](url string, body interface{}, options Options) *T {
	return makeRequest[T](url, http.MethodDelete, body, options)
}

func Option[T any](url string, body interface{}, options Options) *T {
	return makeRequest[T](url, http.MethodOptions, body, options)
}
