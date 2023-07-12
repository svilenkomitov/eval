package eval

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	HeaderContentType   = "Content-Type"
	MimeApplicationJson = "application/json"

	EvaluateEndpoint = "/evaluate"
	ValidateEndpoint = "/validate"
	ErrorsEndpoint   = "/errors"
)

type Handler struct{}

func (h *Handler) Routes(router *chi.Mux) {
	router.MethodFunc(http.MethodPost, EvaluateEndpoint, h.Evaluate)
	router.MethodFunc(http.MethodPost, ValidateEndpoint, h.Validate)
	router.MethodFunc(http.MethodGet, ErrorsEndpoint, h.Errors)
}

func (h *Handler) Evaluate(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusOK, nil)
}

func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusOK, nil)
}

func (h *Handler) Errors(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusOK, nil)
}

func respond(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set(HeaderContentType, MimeApplicationJson)
	w.WriteHeader(code)
	if data != nil {
		resp, _ := json.Marshal(data)
		_, _ = w.Write(resp)
	}
}
