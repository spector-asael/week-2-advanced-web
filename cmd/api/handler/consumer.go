package handler

import (
	"net/http"

	"github.com/spector-asael/week-1/internal/data"
)

func (a *ApplicationDependencies) ShowAllConsumersHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve all consumers from the database
	consumers, err := a.Models.Consumer.GetAll()
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// send the consumers as a JSON response
	a.writeJSON(w, http.StatusOK, envelope{"consumers": consumers}, nil)
}

func (a *ApplicationDependencies) ShowConsumerHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve the consumer ID from the URL
	id := r.PathValue("id")

	// retrieve the consumer from the database
	consumer, err := a.Models.Consumer.GetById(id)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// send the consumer as a JSON response
	a.writeJSON(w, http.StatusOK, envelope{"consumer": consumer}, nil)
}

func (a *ApplicationDependencies) CreateConsumerHandler(w http.ResponseWriter, r *http.Request) {
	// create a struct to hold the incoming consumer data
	var input struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Status string `json:"status"`
	}

	// read the JSON request body into the input struct
	err := a.readJSON(w, r, &input)
	if err != nil {
		a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// create a new consumer instance
	consumer := data.Consumer{
		Name:   input.Name,
		Email:  input.Email,
		Status: input.Status,
	}

	// insert the new consumer into the database
	err = a.Models.Consumer.Insert(&consumer)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// send the newly created consumer as a JSON response
	a.writeJSON(w, http.StatusCreated, envelope{"consumer": consumer}, nil)
}

func (a *ApplicationDependencies) UpdateConsumerHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve the consumer ID from the URL
	id := r.PathValue("id")

	// retrieve the existing consumer
	consumer, err := a.Models.Consumer.GetById(id)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// create a struct to hold the updated consumer data
	var input struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Status string `json:"status"`
	}

	// read the JSON request body into the input struct
	err = a.readJSON(w, r, &input)
	if err != nil {
		a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// update the consumer instance
	consumer.Name = input.Name
	consumer.Email = input.Email
	consumer.Status = input.Status

	// update the consumer in the database
	err = a.Models.Consumer.UpdateById(consumer)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// send the updated consumer as a JSON response
	a.writeJSON(w, http.StatusOK, envelope{"consumer": consumer}, nil)
}
