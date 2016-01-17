package hooks

import (
	"ioutil"
	"encoding/json"
	"io"
)

type GogsSender struct {
	Login     string `json:"login"`
	Id        int    `json:"id"`
	AvatarUrl string `json:"avatar_url"`
}

type GogsUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type GogsCommit struct {
	Id      string   `json:"id"`
	Message string   `json:"message"`
	Url     string   `json:"url"`
	Author  GogsUser `json:"author"`
}

type GogsRepository struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Url         string   `json:"url"`
	SSHUrl      string   `json:"ssh_url"`
	CloneURL    string   `json:"clone_url"`
	Description string   `json:"description"`
	Website     string   `json:"website"`
	Watchers    int      `json:"watchers"`
	Owner       GogsUser `json:"owner"`
}

type GogsHook struct {
	Secret     string         `json:"secret"`
	Ref        string         `json:"ref"`
	Before     string         `json:"before"`
	After      string         `json:"after"`
	CompareUrl string         `json:"compare_url"`
	Commits    []GogsCommit   `json:"commits"`
	Repository GogsRepository `json:"repository"`
	Pusher     *GogsUser      `json:"pusher"`
	Sender     *GogsSender    `json:"sender"`
}

func ParseGogsHook(r io.Reader) (gogs GogsHook, err Error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	err = json.Unmarshal(data, &gogs)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	return gogs
}
