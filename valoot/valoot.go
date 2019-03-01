package valoot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	// APIBackend is a constant representing the API service backend.
	APIBackend SupportedBackend = "api"

	defaultHTTPTimeout = 10 * time.Second
)

var (
	apiUrl   string
	debug    bool
	backends ValootBackend

	httpClient = &http.Client{Timeout: defaultHTTPTimeout}
)

type Backend interface {
	Call(method, path, accessToken string, form *url.Values, content interface{}, v interface{}) error
}

type BackendConfiguration struct {
	Type       SupportedBackend
	URL        string
	HTTPClient *http.Client
}

type SupportedBackend string

type ValootBackend struct {
	API Backend
	mu  sync.RWMutex
}

func (s *BackendConfiguration) NewRequest(method, path, accessToken string, body io.Reader) (*http.Request, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	path = s.URL + path

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		log.Printf("Cannot create valoot request: %v\n", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if accessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	return req, nil
}

func (s *BackendConfiguration) Do(req *http.Request, v interface{}) error {
	log.Printf("Requesting %v %v%v\n", req.Method, req.URL.Host, req.URL.Path)
	start := time.Now()

	res, err := s.HTTPClient.Do(req)

	if debug {
		log.Printf("Completed in %v\n", time.Since(start))
	}

	if err != nil {
		log.Printf("Request to valoot failed: %v\n", err)
		return err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Cannot parse valoot response: %v\n", err)
		return err
	}

	log.Printf("Valoot resp body: %s\n", string(resBody))

	// parses error response if status code is not 2xx
	if res.StatusCode < 200 || res.StatusCode >= 400 {
		var err ErrorResp
		json.Unmarshal(resBody, &err)
		return err
	}

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
func (s BackendConfiguration) Call(method, path, accessToken string, form *url.Values, content interface{}, v interface{}) error {
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

			fmt.Printf("body: %s\n", string(encoded))
			body = bytes.NewBuffer(encoded)
		}
	}

	req, err := s.NewRequest(method, path, accessToken, body)
	if err != nil {
		return err
	}

	if err := s.Do(req, v); err != nil {
		return err
	}

	return nil
}

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

func SetDebug(value bool) {
	debug = value
}

func Setup(url string) {
	apiUrl = url
}

func SetClientTimeout(t time.Duration) {
	httpClient.Timeout = t
}
