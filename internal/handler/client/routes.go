package client

import (
	"encoding/json"
	"errors"
	"github.com/dwnGnL/ddos-pow/config"
	"github.com/dwnGnL/ddos-pow/internal/application"
	"github.com/dwnGnL/ddos-pow/lib/goerrors"
	"github.com/dwnGnL/ddos-pow/lib/pow"
	challengeResp "github.com/dwnGnL/ddos-pow/lib/protocol/challenge-resp"
	"net/http"
)

type Handler struct {
	conf *config.Config
}

func newHandler(cfg *config.Config) *Handler {
	return &Handler{conf: cfg}
}

// Response general response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	SUCCESS = "success" // when we successfully get a quote
	ERROR   = "error"   // when there is server error
	FAIL    = "fail"    // when user has sent invalid data
)

// RequestChallenge to get challenge to solve from server
func (h Handler) RequestChallenge(w http.ResponseWriter, r *http.Request) {
	app, err := application.GetAppFromRequest(r)
	if err != nil {
		respond(w, http.StatusBadRequest, nil, err)
		return
	}

	clientService := app.GetClient()
	hashCashData, err := clientService.RequestChallenge()

	if err != nil {
		respond(w, http.StatusInternalServerError, nil, err)
		return
	}

	respond(w, http.StatusOK, hashCashData, nil)
}

// RequestResource to get quote from server
func (h Handler) RequestResource(w http.ResponseWriter, r *http.Request) {
	app, err := application.GetAppFromRequest(r)
	if err != nil {
		respond(w, http.StatusBadRequest, nil, err)
		return
	}

	hashCashData := pow.HashcashData{}

	// unmarshall the challenge
	err = json.NewDecoder(r.Body).Decode(&hashCashData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	goerrors.Log().Println("body hashcashData", hashCashData)

	clientService := app.GetClient()
	msg, err := clientService.RequestResource(hashCashData)

	// server error
	if err != nil {
		respond(w, http.StatusInternalServerError, nil, err)
		return
	}
	// invalid challenge data has been sent to server
	if msg.Header == challengeResp.FAIL {
		respond(w, http.StatusBadRequest, nil, errors.New(msg.Payload))
		return
	}

	// respond the quote
	respond(w, http.StatusOK, msg.Payload, nil)
}

func respond(w http.ResponseWriter, status int, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")

	var response Response

	w.WriteHeader(status)

	switch {
	case status >= 200 && status <= 299:
		response.Status = SUCCESS
		response.Data = data
	case status >= 400 && status < 499:
		response.Status = FAIL
		response.Message = err.Error()
	default:
		goerrors.Log().Println(err)
		response.Status = ERROR
		response.Message = "internal server error"
	}

	resp, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
