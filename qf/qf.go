package qf

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type SupportedBackend string

const (
	// APIBackend is a constant representing the API service backend.
	APIBackend SupportedBackend = "api"

	defaultHTTPTimeout = 10 * time.Second
)

// AppCode is the qf code used globally in the binding.
var AppCode string

var (
	apiUrl     string
	backends   Backends
	httpClient = &http.Client{Timeout: defaultHTTPTimeout}
)

type Backend interface {
	Call(method, path, qfAppCode, qfSign string, form *url.Values, content interface{}, v interface{}) error
}

type BackendConfiguration struct {
	Type       SupportedBackend
	URL        string
	HTTPClient *http.Client
}

type Backends struct {
	API Backend
	mu  sync.RWMutex
}

func (s *BackendConfiguration) NewRequest(method, path, qfAppCode, qfSign string, body io.Reader) (*http.Request, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = s.URL + path

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		log.Printf("Cannot create qf request: %v\n", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-QF-APPCODE", qfAppCode)
	req.Header.Add("X-QF-SIGN", qfSign)

	return req, nil
}

func (s *BackendConfiguration) Do(req *http.Request, v interface{}) error {
	log.Printf("Requesting %v %v%v\n", req.Method, req.URL.Host, req.URL.Path)

	res, err := s.HTTPClient.Do(req)

	if err != nil {
		log.Printf("Failed to do http request: %v\n", err)
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("Failed to send request: %v\n", err)
		return err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read res.body: %v\n", err)
		return err
	}

	log.Printf("QF resp body: %s\n", string(resBody))

	if v != nil {
		return json.Unmarshal(resBody, v)
	}

	return nil
}

/**
 * @param method Either GET, POST, PUT, DELETE
 * @param path URL path
 * @param form Query string for GET method only
 * @param content Interface{} of JSON object
 * @param v Any response object and fill after call success
 * @example
 *
 * // GET /some_resource?a=1
 * u := &url.Values{}l
 * u.Add("a", "1")
 * obj := RespObj{}
 * s.Call("GET", "/some_url", "xxx", u, nil, obj)
 *
 * // POST /some_resource
 * c := currency.Content{}
 * c["data"] = "John"
 * obj := RespObj{}
 * s.Call("POST", "/some_resource", "xxx", "1234567890", nil, &c, obj)
 */
func (s *BackendConfiguration) Call(method, path, qfAppCode, qfSign string, form *url.Values, content interface{}, v interface{}) error {
	var body io.Reader

	method = strings.ToUpper(method)
	if method == "GET" {
		if form != nil && len(*form) > 0 {
			path += "?" + form.Encode()
		}
	} else {
		// POST, PUT, DELETE
		if content != nil {
			encoded, _ := json.Marshal(content)

			log.Printf("encoded body: %s\n", string(encoded))
			body = bytes.NewBuffer(encoded)
		}
	}

	req, err := s.NewRequest(method, path, qfAppCode, qfSign, body)
	if err != nil {
		return err
	}

	if err := s.Do(req, v); err != nil {
		return err
	}

	return nil
}

// Concurrency-safe
func GetBackend(backendType SupportedBackend) Backend {
	var backend Backend

	backends.mu.RLock()

	switch backendType {
	case APIBackend:
		if backends.API != nil {
			backend = backends.API
			backends.mu.RUnlock()

			return backend
		}
	}

	backends.mu.RUnlock()

	// acquire an exclusive lock
	backends.mu.Lock()

	// must check for nil
	switch backendType {
	case APIBackend:
		if backends.API == nil {
			backends.API = &BackendConfiguration{backendType, apiUrl, httpClient}
		}
		backend = backends.API
	}

	backends.mu.Unlock()
	return backend
}

func SetUrl(url string) {
	apiUrl = url
}

func SetClientTimeout(t time.Duration) {
	httpClient.Timeout = t
}
