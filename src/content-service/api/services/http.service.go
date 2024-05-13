package services

import "net/http"

type HttpService interface {
	http.Handler
}

type httpService struct {
	mux *http.ServeMux
}

func (h httpService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func NewHttpService() (HttpService, error) {
	return httpService{mux: http.NewServeMux()}, nil
}
