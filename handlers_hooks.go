package main

import (
	"net/http"
	"strings"

	. "github.com/jakecoffman/gorunner/service"
	"github.com/nu7hatch/gouuid"
	"github.com/jakecoffman/gorunner/hooks"
)

type msi map[string]interface{}

type GenericHook struct {
	RepoName string
	RepoURL string
	LastCommit string
	Branch string
}

func runHook(c context, gh GenericHook) (int, interface{}) {
	var envs []Environment
	envs = append(envs, Environment{Name: "branch", Id: "CI_BRANCH", Value: gh.Branch})
	envs = append(envs, Environment{Name: "repourl", Id: "CI_REPOURL", Value: gh.RepoURL})
	envs = append(envs, Environment{Name: "reponame", Id: "CI_REPO", Value: gh.RepoName})
	envs = append(envs, Environment{Name: "commit", Id: "CI_COMMIT", Value: gh.LastCommit})

	var ids []string
	jobs := c.JobList().GetJobsWithTrigger(gh.RepoName)
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

func hookGogs(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	msg, err := hooks.ParseGogsHook(r.Body)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	reponame := msg.Repository.Name
	repourl := msg.Repository.CloneURL
	last_commit := msg.After
	ref := strings.Split(msg.Ref, "/")

	// Gogs is stupid sometimes,
	// giving ref="master" instead of ref="refs/heads/master"
	branch := ref[0]
	if len(ref) == 3 {
		branch = ref[2]
	}

	gh := GenericHook{ RepoName: reponame, RepoURL: repourl, LastCommit: last_commit, Branch: branch }
	return runHook(c, gh)
}

func hookBitbucket(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	msg, err := hooks.ParseBitbucketHook(r.Body)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	reponame := msg.Repository.Name
	last_commit := msg.Push.Changes[0].New.Target.Hash
	branch := msg.Push.Changes[0].New.Name
	repourl := "git@bitbucket.org:" + msg.Repository.FullName + ".git"

	gh := GenericHook{ RepoName: reponame, RepoURL: repourl, LastCommit: last_commit, Branch: branch }
	return runHook(c, gh)
}
