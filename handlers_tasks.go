package main

import (
	"net/http"

	"github.com/gorilla/mux"
	. "github.com/jakecoffman/gorunner/service"
)

// Tasks

func listTasks(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	return http.StatusOK, c.TaskList().Dump()
}

func addTask(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	payload := unmarshal(r.Body, "name", w)

	err := c.TaskList().Append(Task{payload["name"], ""})
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}
	return http.StatusCreated, nothing
}

func getTask(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	task, err := c.TaskList().Get(vars["task"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	return http.StatusOK, task
}

func updateTask(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	task, err := c.TaskList().Get(vars["task"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	payload := unmarshal(r.Body, "script", w)
	t := task.(Task)
	t.Script = payload["script"]
	c.TaskList().Update(t)
	return http.StatusOK, nothing
}

func deleteTask(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	task, err := c.TaskList().Get(vars["task"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	c.TaskList().Delete(task.ID())
	return http.StatusOK, nothing
}

func listJobsForTask(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	jobs := c.JobList().GetJobsWithTask(vars["task"])
	return http.StatusOK, jobs
}
