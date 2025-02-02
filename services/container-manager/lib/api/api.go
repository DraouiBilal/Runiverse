package api

import (
	"bufio"
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
	Debug   bool
	Socket  string
	Follow  bool
}

func createRequest(url string, method string, body interface{}, options Options) *http.Request {
	payload, err := json.Marshal(body)

	if err != nil {
		log.Fatal(err)
	}

	if options.Debug {
		log.Println(bytes.NewBuffer(payload))
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

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

	if options.Debug {
		requestDump, err := httputil.DumpRequest(req, true)
		responseDump, reerr := httputil.DumpResponse(res, true)
		if err != nil {
			log.Fatal(err)
		}

		if reerr != nil {
			log.Fatal(err)
		}
		log.Println(string(requestDump))
		log.Println(string(responseDump))
	}

	if err != nil || res.StatusCode >= 400 {
		log.Println("Request failed or Status code >= 400")
		responseDump, err := httputil.DumpResponse(res, true)
		if err != nil {
			panic(err)
		}
		panic(string(responseDump))
	}

	var result T

	switch any(result).(type) {
	case string:
		resText := ""
        var err interface{}

		if options.Follow {
			scanner := bufio.NewScanner(res.Body)

			for scanner.Scan() {
				resText += scanner.Text() + "\n"
			}

            err = scanner.Err()
		} else {
            var resBytes []byte
			resBytes, err = io.ReadAll(res.Body)
            resText = string(resBytes)
		}

		if err != nil {
			log.Fatal(err)
		}
		strResult := string(resText)
		return any(&strResult).(*T)
	}

	if reflect.TypeOf((*T)(nil)).Elem().Kind() == reflect.Interface {
		return nil
	}

	responseJson := new(T)

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&responseJson)

	if err != nil {
		log.Println("JSON decoder failed")
		log.Fatal(err)
	}

	defer res.Body.Close()

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
