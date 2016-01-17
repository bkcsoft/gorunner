package hooks

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type BitBucketLink struct {
	HREF string `json:"href"`
}

type BitbucketActor struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Type        string `json:"user"`
	UUID        string `json:"uuid"`
	Links       struct {
		Self   BitBucketLink `json:"self"`
		HTML   BitBucketLink `json:"html"`
		Avatar BitBucketLink `json:"avatar"`
	} `json:"links"`
}

type BitbucketCommit struct {
	Hash    string      `json:"hash"`
	Parent  interface{} `json:"parent"`
	Date    string      `json:"date"`
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Links   struct {
		Self BitBucketLink `json:"self"`
		HTML BitBucketLink `json:"html"`
	} `json:"links"`
	Author struct {
		Raw  string         `json:"raw"`
		User BitbucketActor `json:"user"`
	} `json:"author"`
}

type BitbucketChange struct {
	Type       string          `json:"type"`
	Name       string          `json:"name"`
	Repository interface{}     `json:"repository"`
	Target     BitbucketCommit `json:"target"`
	Links      struct {
		Commits BitBucketLink `json:"commits"`
		HTML    BitBucketLink `json:"html"`
		Self    BitBucketLink `json:"self"`
	} `json:"links"`
}

type BitbucketChangeSet struct {
	Forced    bool            `json:"forced"`
	Created   bool            `json:"created"`
	Truncated bool            `json:"truncated"`
	Closed    bool            `json:"closed"`
	Old       BitbucketChange `json:"old"`
	New       BitbucketChange `json:"new"`
	Links     struct {
		Diff    BitBucketLink `json:"diff"`
		HTML    BitBucketLink `json:"html"`
		Commits BitBucketLink `json:"commits"`
	} `json:"links"`
}

type BitbucketPush struct {
	Changes []BitbucketChangeSet `json:"changes"`
}

type BitbucketRepository struct {
	Website   string         `json:"website"`
	SCM       string         `json:"scm"`
	Name      string         `json:"Fuckemon"`
	FullName  string         `json:"full_name"`
	Owner     BitbucketActor `json:"owner"`
	Type      string         `json:"type"`
	IsPrivate bool           `json:"is_private"`
	UUID      string         `json:"uuid"`
	Links     struct {
		Self   BitBucketLink `json:"self"`
		HTML   BitBucketLink `json:"html"`
		Avatar BitBucketLink `json:"avatar"`
	} `json:"links"`
}

type BitbucketHook struct {
	Push       BitbucketPush       `json:"push"`
	Repository BitbucketRepository `json:"repository"`
	Actor      BitbucketActor      `json:"actor"`
}

func ParseBitbucketHook(r io.Reader) (bitbucket *BitbucketHook, err error) {
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
