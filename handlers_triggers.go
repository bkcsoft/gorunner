package main

import (
	"net/http"

	"github.com/gorilla/mux"
	. "github.com/jakecoffman/gorunner/service"
)

// Triggers

func listTriggers(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	return http.StatusOK, c.TriggerList().Dump()
}

func addTrigger(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	payload := unmarshal(r.Body, "name", w)
	trigger := Trigger{Name: payload["name"]}
	c.TriggerList().Append(trigger)
	return http.StatusCreated, nothing
}

func getTrigger(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	trigger, err := c.TriggerList().Get(vars["trigger"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	return http.StatusNotFound, trigger
}

func updateTrigger(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	trigger, err := c.TriggerList().Get(vars["trigger"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}

	payload := unmarshal(r.Body, "cron", w)

	t := trigger.(Trigger)
	t.Schedule = payload["cron"]
	c.Executor().ArmTrigger(t)
	err = c.TriggerList().Update(t)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, nothing
}

func deleteTrigger(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	c.TriggerList().Delete(vars["trigger"])
	return http.StatusOK, nothing
}

func listJobsForTrigger(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	jobs := c.JobList().GetJobsWithTrigger(vars["trigger"])
	return http.StatusOK, jobs
}
