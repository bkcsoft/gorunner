package main

import (
	"log"
	"net/http"
	"strings"

	. "github.com/jakecoffman/gorunner/service"
	"github.com/nu7hatch/gouuid"
	"github.com/jakecoffman/gorunner/hooks"
)

type msi map[string]interface{}


func hookGogs(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	var (
		reponame    string
		repourl     string
		last_commit string
		ref         []string
	)
	msg := unmarshalAll(r.Body, w)
	faster, ok := msg.(gogsHook)
	if ok {
		reponame = faster.Repository.Name
		repourl = faster.Repository.CloneURL
		last_commit = faster.After
		ref = strings.Split(faster.Ref, "/")
	} else {
		slower := msg.(map[string]interface{})
		log.Println("Can't cast hook to struct, falling back to mapping manually...")
		reponame = slower["repository"].(map[string]interface{})["name"].(string)
		repourl = slower["repository"].(map[string]interface{})["clone_url"].(string)
		last_commit = slower["after"].(string)
		ref = strings.Split(slower["ref"].(string), "/")
		/*
			// Not in Gogs yet...
			if msg["created"] != true {
				return http.StatusOK, nothing
			}
		*/
	}

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

// {name: .repository.name, commit: .push.changes[0].new.target.bash, branch: .push.changes[0].new.name}
func hookBitbucket(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	msg := unmarshalAll(r.Body, w).(map[string]interface{})

	reponame := string.ToLower(msg["repository"].(msi)["name"])
	repourl := "magic..."
	last_commit := msg["push"].(msi)["changes"].([]interface{})[0].(msi)["new"].(msi)["target"].(msi)["name"].(string)
	branch := msg["push"].(msi)["changes"].([]interface{})[0].(msi)["new"].(msi)["name"].(string)

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
