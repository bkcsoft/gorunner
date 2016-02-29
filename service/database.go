package service

import (
	"io/ioutil"
	"os"
	"sort"
)

const (
	jobsFile     = "db/jobs.json"
	runsFile     = "db/runs.json"
	tasksFile    = "db/tasks.json"
	triggersFile = "db/triggers.json"
)

type ListWriter func([]byte, string)
type ListReader func(string) []byte

func writeFile(bytes []byte, filePath string) {
	err := ioutil.WriteFile(filePath, bytes, 0644)
	if err != nil {
		panic(err)
	}
}

func readFile(filePath string) []byte {
	_, err := os.Stat(filePath)
	if err != nil {
		println("Couldn't read file, creating fresh:", filePath)
		writeFile([]byte("[]"), filePath)
	}

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return bytes
}

type Reverse struct {
	sort.Interface
}

func (r Reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}
