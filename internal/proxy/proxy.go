package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Handler struct {
	reverseProxy *httputil.ReverseProxy
}

func NewHandler(backendURL string) (*Handler, error) {
	target, err := url.Parse(backendURL)
	if err != nil {
		return nil, err
	}

	rp := httputil.NewSingleHostReverseProxy(target)
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	return &Handler{reverseProxy: rp}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.reverseProxy.ServeHTTP(w, r)
}
