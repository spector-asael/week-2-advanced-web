package handler

import (
	"net/http"
	"time"

	"github.com/spector-asael/week-1/internal/data"
)

func (a *ApplicationDependencies) ShowAllAPIKeysHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve all API keys from the database
	apiKeys, err := a.Models.API_Keys.GetAll()
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// send the API keys as a JSON response
	a.writeJSON(w, http.StatusOK, envelope{"api_keys": apiKeys}, nil)
}

func (a *ApplicationDependencies) ShowAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve the API key ID from the URL
	id := r.PathValue("id")

	// retrieve the API key from the database
	apiKey, err := a.Models.API_Keys.GetById(id)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// send the API key as a JSON response
	a.writeJSON(w, http.StatusOK, envelope{"api_key": apiKey}, nil)
}

func (a *ApplicationDependencies) CreateAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	// create a struct to hold the incoming API key data
	var input struct {
		ConsumerID string     `json:"consumer_id"`
		KeyHash    string     `json:"key_hash"`
		KeyPrefix  string     `json:"key_prefix"`
		Status     string     `json:"status"`
		LastUsedAt *time.Time `json:"last_used_at"`
		ExpiresAt  *time.Time `json:"expires_at"`
	}

	// read the JSON request body into the input struct
	err := a.readJSON(w, r, &input)
	if err != nil {
		a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// create a new API key instance
	apiKey := data.API_Keys{
		ConsumerID: input.ConsumerID,
		KeyHash:    input.KeyHash,
		KeyPrefix:  input.KeyPrefix,
		Status:     input.Status,
		LastUsedAt: input.LastUsedAt,
		ExpiresAt:  input.ExpiresAt,
	}

	// insert the new API key into the database
	err = a.Models.API_Keys.Insert(&apiKey)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// send the newly created API key as a JSON response
	a.writeJSON(w, http.StatusCreated, envelope{"api_key": apiKey}, nil)
}
