package eval

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	HeaderContentType   = "Content-Type"
	MimeApplicationJson = "application/json"

	EvaluateEndpoint = "/evaluate"
	ValidateEndpoint = "/validate"
	ErrorsEndpoint   = "/errors"
)

type EvaluationRequest struct {
	Expression string `json:"expression"`
}

type EvaluationResponse struct {
	Result int `json:"result"`
}

type ValidationRequest struct {
	Expression string `json:"expression"`
}

type ValidationResponse struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason,omitempty"`
}

type ErrorResponse struct {
	Expression string `json:"expression"`
	Endpoint   string `json:"endpoint"`
	Frequency  int    `json:"frequency"`
	Type       string `json:"type"`
}

type Handler struct {
	Service Service
}

func (h *Handler) Routes(router *chi.Mux) {
	router.MethodFunc(http.MethodPost, EvaluateEndpoint, h.Evaluate)
	router.MethodFunc(http.MethodPost, ValidateEndpoint, h.Validate)
	router.MethodFunc(http.MethodGet, ErrorsEndpoint, h.Errors)
}

func (h *Handler) Evaluate(w http.ResponseWriter, r *http.Request) {
	var requestBody EvaluationRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Errorf("invalid request body: %v", err)
		respondError(w, &ResponseError{
			Kind:   ValidationError,
			Reason: "invalid request body",
			Code:   http.StatusBadRequest,
		})
		return
	}

	result, err := h.Service.Evaluate(requestBody.Expression)
	if err != nil {
		respondError(w, err)
		return
	}

	respond(w, http.StatusOK, EvaluationResponse{
		Result: result,
	})
}

func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {
	var requestBody ValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Errorf("invalid request body: %v", err)
		respondError(w, &ResponseError{
			Kind:   ValidationError,
			Reason: "invalid request body",
			Code:   http.StatusBadRequest,
		})
		return
	}

	if err := h.Service.Validate(requestBody.Expression); err != nil {
		respond(w, http.StatusOK, ValidationResponse{
			Valid:  false,
			Reason: err.Error(),
		})
		return
	}

	respond(w, http.StatusOK, ValidationResponse{
		Valid: true,
	})
}

func (h *Handler) Errors(w http.ResponseWriter, r *http.Request) {
	errors, err := h.Service.Errors()
	if err != nil {
		respondError(w, err)
		return
	}
	respond(w, http.StatusOK, errors)
}

func respondError(w http.ResponseWriter, err *ResponseError) {
	respond(w, err.Code, err)
}

func respond(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set(HeaderContentType, MimeApplicationJson)
	w.WriteHeader(code)
	if data != nil {
		resp, _ := json.Marshal(data)
		_, _ = w.Write(resp)
	}
}
