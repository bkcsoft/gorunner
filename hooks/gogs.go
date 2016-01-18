package hooks

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type gogsSender struct {
	Login     string `json:"login"`
	Id        int    `json:"id"`
	AvatarUrl string `json:"avatar_url"`
}

type gogsUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type gogsCommit struct {
	Id      string   `json:"id"`
	Message string   `json:"message"`
	Url     string   `json:"url"`
	Author  gogsUser `json:"author"`
}

type gogsRepository struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Url         string   `json:"url"`
	SSHUrl      string   `json:"ssh_url"`
	CloneURL    string   `json:"clone_url"`
	Description string   `json:"description"`
	Website     string   `json:"website"`
	Watchers    int      `json:"watchers"`
	Owner       gogsUser `json:"owner"`
}

type gogsHook struct {
	Secret     string         `json:"secret"`
	Ref        string         `json:"ref"`
	Before     string         `json:"before"`
	After      string         `json:"after"`
	CompareUrl string         `json:"compare_url"`
	Commits    []gogsCommit   `json:"commits"`
	Repository gogsRepository `json:"repository"`
	Pusher     *gogsUser      `json:"pusher"`
	Sender     *gogsSender    `json:"sender"`
}

func ParseGogsHook(r io.Reader) (gogs *gogsHook, err error) {
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

	return gogs, nil
}
