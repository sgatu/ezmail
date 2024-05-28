package handlers

import "net/http"

type helloHandler struct{}

func (h *helloHandler) hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
