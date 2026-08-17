package handler

import (
	"net/http"

	"github.com/spector-asael/week-1/internal/data"
)

func (a *ApplicationDependencies) ShowAllJobsHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve all jobs from the database
	jobs, err := a.Models.Jobs.GetAll()
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// send the jobs as a JSON response
	a.writeJSON(w, http.StatusOK, envelope{"jobs": jobs}, nil)
}

func (a *ApplicationDependencies) CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	// create a struct to hold the incoming job data
	var input struct {
		ConsumerID string `json:"consumer_id"`
		JobType    string `json:"job_type"`
		Status     string `json:"status"`
		Payload    string `json:"payload"`
		Result     string `json:"result"`
	}

	// read the JSON request body into the input struct
	err := a.readJSON(w, r, &input)
	if err != nil {
		a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// create a new job instance
	job := data.Jobs{
		ConsumerID: input.ConsumerID,
		JobType:    input.JobType,
		Status:     input.Status,
		Payload:    input.Payload,
		Result:     &input.Result,
	}

	// insert the new job into the database
	err = a.Models.Jobs.Insert(&job)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// send a JSON response with the newly created job
	a.writeJSON(w, http.StatusCreated, envelope{"job": job}, nil)
}
