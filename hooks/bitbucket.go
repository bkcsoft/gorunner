package hooks

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type bitbucketLink struct {
	HREF string `json:"href"`
}

type bitbucketActor struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Type        string `json:"user"`
	UUID        string `json:"uuid"`
	Links       struct {
		Self   bitbucketLink `json:"self"`
		HTML   bitbucketLink `json:"html"`
		Avatar bitbucketLink `json:"avatar"`
	} `json:"links"`
}

type bitbucketCommit struct {
	Hash    string      `json:"hash"`
	Parent  interface{} `json:"parent"`
	Date    string      `json:"date"`
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Links   struct {
		Self bitbucketLink `json:"self"`
		HTML bitbucketLink `json:"html"`
	} `json:"links"`
	Author struct {
		Raw  string         `json:"raw"`
		User bitbucketActor `json:"user"`
	} `json:"author"`
}

type bitbucketChange struct {
	Type       string          `json:"type"`
	Name       string          `json:"name"`
	Repository interface{}     `json:"repository"`
	Target     bitbucketCommit `json:"target"`
	Links      struct {
		Commits bitbucketLink `json:"commits"`
		HTML    bitbucketLink `json:"html"`
		Self    bitbucketLink `json:"self"`
	} `json:"links"`
}

type bitbucketChangeSet struct {
	Forced    bool            `json:"forced"`
	Created   bool            `json:"created"`
	Truncated bool            `json:"truncated"`
	Closed    bool            `json:"closed"`
	Old       bitbucketChange `json:"old"`
	New       bitbucketChange `json:"new"`
	Links     struct {
		Diff    bitbucketLink `json:"diff"`
		HTML    bitbucketLink `json:"html"`
		Commits bitbucketLink `json:"commits"`
	} `json:"links"`
}

type bitbucketPush struct {
	Changes []bitbucketChangeSet `json:"changes"`
}

type bitbucketRepository struct {
	Website   string         `json:"website"`
	SCM       string         `json:"scm"`
	Name      string         `json:"name"`
	FullName  string         `json:"full_name"`
	Owner     bitbucketActor `json:"owner"`
	Type      string         `json:"type"`
	IsPrivate bool           `json:"is_private"`
	UUID      string         `json:"uuid"`
	Links     struct {
		Self   bitbucketLink `json:"self"`
		HTML   bitbucketLink `json:"html"`
		Avatar bitbucketLink `json:"avatar"`
	} `json:"links"`
}

type bitbucketHook struct {
	Push       bitbucketPush       `json:"push"`
	Repository bitbucketRepository `json:"repository"`
	Actor      bitbucketActor      `json:"actor"`
}

func ParseBitbucketHook(r io.Reader) (bitbucket *bitbucketHook, err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	err = json.Unmarshal(data, &bitbucket)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	return bitbucket, nil
}
