package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	. "github.com/jakecoffman/gorunner/service"
)

// Jobs

func listJobs(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	return http.StatusOK, c.JobList().Dump()
}

func addJob(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	payload := unmarshal(r.Body, "name", w)

	err := c.JobList().Append(Job{Name: payload["name"], Status: "New"})
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusCreated, nothing
}

func getJob(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.JobList().Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}

	return http.StatusOK, job
}

func deleteJob(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.JobList().Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}

	err = c.JobList().Delete(job.ID())
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, nothing
}

func addTaskToJob(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.JobList().Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	j := job.(Job)

	payload := unmarshal(r.Body, "task", w)
	j.AppendTask(payload["task"])
	c.JobList().Update(j)

	return http.StatusCreated, nothing
}

func removeTaskFromJob(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.JobList().Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	j := job.(Job)

	taskPosition, err := strconv.Atoi(vars["task"])
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}
	j.DeleteTask(taskPosition)
	c.JobList().Update(j)
	return http.StatusOK, nothing
}

func addTriggerToJob(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.JobList().Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	j := job.(Job)

	payload := unmarshal(r.Body, "trigger", w)

	j.AppendTrigger(payload["trigger"])
	t, err := c.TriggerList().Get(payload["trigger"])
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	c.Executor().ArmTrigger(t.(Trigger))
	c.JobList().Update(j)

	return http.StatusCreated, nothing
}

func removeTriggerFromJob(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.JobList().Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	j := job.(Job)

	t := vars["trigger"]
	j.DeleteTrigger(t)
	c.JobList().Update(j)

	// If Trigger is no longer attached to any Jobs, remove it from Cron to save cycles
	jobs := c.JobList().GetJobsWithTrigger(t)

	if len(jobs) == 0 {
		c.Executor().DisarmTrigger(t)
	}
	return http.StatusOK, nothing
}
