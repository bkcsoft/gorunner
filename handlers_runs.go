package main

import (
	"net/http"
	"strconv"
	"io/ioutil"
	"encoding/json"

	"github.com/gorilla/mux"
	. "github.com/jakecoffman/gorunner/service"
	"github.com/nu7hatch/gouuid"
)

// Run

func listRuns(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	offset := r.FormValue("offset")
	length := r.FormValue("length")

	if offset == "" {
		offset = "-1"
	}
	if length == "" {
		length = "-1"
	}

	o, err := strconv.Atoi(offset)
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}

	l, err := strconv.Atoi(length)
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}

	return http.StatusOK, c.RunList().GetRecent(o, l)
}

func addRun(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	payload := unmarshal(r.Body, "job", w)

	job, err := c.JobList().Get(payload["job"])
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	j := job.(Job)

	id, err := uuid.NewV4()
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	var tasks []Task
	for _, taskName := range j.Tasks {
		task, err := c.TaskList().Get(taskName)
		if err != nil {
			panic(err)
		}
		t := task.(Task)
		tasks = append(tasks, t)
	}
	err = c.RunList().AddRun(id.String(), j, tasks, "{}")
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusCreated, map[string]string{"uuid": id.String()}
}

func addRun2(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
//	payload := unmarshal(r.Body, "repository", w)
	payload, err := ioutil.ReadAll(r.Body);
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	vars := mux.Vars(r)
	job, err := c.JobList().Get(vars["job"])
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	j := job.(Job)

	id, err := uuid.NewV4()
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	var tasks []Task
	for _, taskName := range j.Tasks {
		task, err := c.TaskList().Get(taskName)
		if err != nil {
			panic(err)
		}
		t := task.(Task)
		tasks = append(tasks, t)
	}
	err = c.RunList().AddRun(id.String(), j, tasks, byteToString(payload))
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusCreated, map[string]string{"uuid": id.String()}

}

func getRun(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	run, err := c.RunList().Get(vars["run"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	return http.StatusOK, run
}

func deleteRuns(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	var runs []string
	err = json.Unmarshal(data, &runs)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	for _, run := range runs {
		err := c.RunList().DeleteRun(run)
		if err != nil {
			return http.StatusNotFound, err.Error()
		}
	}
	return http.StatusOK, nothing
}

func deleteRun(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	err := c.RunList().DeleteRun(vars["run"])
	if(err != nil) {
		return http.StatusNotFound, err.Error()
	}
	return http.StatusOK, nothing
}

