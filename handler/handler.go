// Package handlers provides HTTP request handlers.
package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang/glog"
	"url-shortener/storage"
)

// New returns an http handler for the url shortener.
func New(prefix string, storage storage.Service) http.Handler {
	mux := http.NewServeMux()
	h := handler{prefix, storage}
	mux.HandleFunc("/encode/", responseHandler(h.encode))
	mux.HandleFunc("/", h.redirect)
	mux.HandleFunc("/info/", responseHandler(h.decode))
	return mux
}

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"response"`
}

type handler struct {
	prefix  string
	storage storage.Service
}

func responseHandler(h func(io.Writer, *http.Request) (interface{}, int, error)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		data, status, err := h(w, r)
		if err != nil {
			data = err.Error()
		}
		fmt.Println("data", data)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(response{Data: data, Success: err == nil})
		if err != nil {
			glog.V(0).Infof("could not encode response to output: %v", err)
		}
	}
}

func (h handler) encode(w io.Writer, r *http.Request) (interface{}, int, error) {

	fmt.Println("r.Method====", r)
	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("method %s not allowed", r.Method)
	}

	var input struct {
		URL        string
		ClientType string
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Unable to decode JSON request body: %v", err)
	}
	//添加判断客户端

	fmt.Println("User-Agent", r.Header.Get("User-Agent"))
	var clientType string
	url := strings.TrimSpace(input.URL)
	if url == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("URL is empty")
	}

	clientType = r.Header.Get("User-Agent")

	short, err := h.storage.Save(url, clientType)
	fmt.Println("short", short)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Could not store in database: %v", err)
	}

	return h.prefix + short, http.StatusCreated, nil
}

func (h handler) decode(w io.Writer, r *http.Request) (interface{}, int, error) {
	if r.Method != http.MethodGet {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("Method %s not allowed", r.Method)
	}

	code := r.URL.Path[len("/info/"):]

	model, err := h.storage.LoadInfo(code)
	if err != nil {
		fmt.Println(err)
		return nil, http.StatusNotFound, fmt.Errorf("URL not found")
	}

	return model, http.StatusOK, nil
}

func (h handler) redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	code := r.URL.Path[len("/"):]

	url, err := h.storage.Load(code)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("URL Not Found"))
		fmt.Println("err", err)
		return
	}
	http.Redirect(w, r, string(url), http.StatusMovedPermanently)
}
