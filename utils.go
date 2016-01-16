package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func byteToString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

func marshal(item interface{}, w http.ResponseWriter) {
	bytes, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(bytes)
}

func unmarshal(r io.Reader, k string, w http.ResponseWriter) (payload map[string]string) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(data, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if payload[k] == "" {
		http.Error(w, "Please provide a '"+k+"'", http.StatusBadRequest)
		return
	}

	return
}

func unmarshalAll(r io.Reader, w http.ResponseWriter) (payload interface{}) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(data, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	return
}
