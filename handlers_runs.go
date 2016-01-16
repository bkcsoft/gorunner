package main

import (
	"net/http"
	"strconv"
	"strings"

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
	err = c.RunList().AddRun(id.String(), j, tasks, nil)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusCreated, map[string]string{"uuid": id.String()}
}

func hookGogs(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	msg := unmarshalAll(r.Body, w).(map[string]interface{})

/*
	// Not in Gogs yet...
	if msg["created"] != true {
		return http.StatusOK, nothing
	}
*/

	reponame := msg["repository"].(map[string]interface{})["name"].(string)
	repourl := msg["repository"].(map[string]interface{})["clone_url"].(string)
	last_commit := msg["after"].(string)
	ref := strings.Split(msg["ref"].(string), "/")
	// Gogs is stupid somethimes giving ref="master" instead of ref="refs/heads/master"
	branch := ref[0]
	if len(ref) == 3 {
		branch = ref[2]
	}

	var envs []Environment
	envs = append(envs, Environment{Name: "branch", Id: "CI_BRANCH", Value: branch})
	envs = append(envs, Environment{Name: "repourl", Id: "CI_REPOURL", Value: repourl})
	envs = append(envs, Environment{Name: "reponame", Id: "CI_REPO", Value: reponame})
	envs = append(envs, Environment{Name: "commit", Id: "CI_COMMIT", Value: last_commit})


	var ids []string
	jobs := c.JobList().GetJobsWithTrigger(reponame)
	for _, j := range jobs {
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
		err = c.RunList().AddRun(id.String(), j, tasks, envs)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}
		ids = append(ids, id.String())

	}

	return http.StatusCreated, ids
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
	runs := unmarshalAll(r.Body, w)
	for _, run := range runs.([]string) {
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
