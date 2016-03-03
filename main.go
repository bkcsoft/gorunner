package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/service"

	"gopkg.in/ini.v1"
)

// AppSettings stores Application Settings...
type AppSettings struct {
	AppName  string `ini:"APP_NAME"`
	HTTPPort int    `ini:"HTTP_PORT"`
	ServeIP  string `ini:"SERVE_IP"`
	LogFile  string `ini:"LOG_FILE"`
}

var (
	// Settings \o/
	Settings AppSettings
)

func loadSettings(filename string) {
	err := ini.MapTo(&Settings, filename)
	if err != nil {
		log.Println("Failed to load `" + filename + "`. Using Defaults instead")
		log.Println(err)
		Settings.HTTPPort = 8090
		Settings.ServeIP = "0.0.0.0"
		Settings.LogFile = ""
	}
	// Set defaults if omited...
	if Settings.HTTPPort == 0 {
		Settings.HTTPPort = 8090
	}
	if Settings.ServeIP == "" {
		Settings.ServeIP = "0.0.0.0"
	}
}

var routes = []struct {
	route   string
	handler func(context, http.ResponseWriter, *http.Request) (int, interface{})
	method  string
}{
	{"/jobs", listJobs, "GET"},
	{"/jobs", addJob, "POST"},
	{"/jobs/{job}", getJob, "GET"},
	{"/jobs/{job}", deleteJob, "DELETE"},
	{"/jobs/{job}/tasks", addTaskToJob, "POST"},
	{"/jobs/{job}/tasks/{task}", removeTaskFromJob, "DELETE"},
	{"/jobs/{job}/triggers", addTriggerToJob, "POST"},
	{"/jobs/{job}/triggers/{trigger}", removeTriggerFromJob, "DELETE"},

	{"/tasks", listTasks, "GET"},
	{"/tasks", addTask, "POST"},
	{"/tasks/{task}", getTask, "GET"},
	{"/tasks/{task}", updateTask, "PUT"},
	{"/tasks/{task}", deleteTask, "DELETE"},
	{"/tasks/{task}/jobs", listJobsForTask, "GET"},

	{"/hooks/gogs", hookGogs, "POST"},
	{"/hooks/bitbucket", hookBitbucket, "POST"},

	{"/runs", listRuns, "GET"},
	{"/runs", addRun, "POST"},
	{"/runs", deleteRuns, "DELETE"},
	{"/runs/{run}", getRun, "GET"},
	{"/runs/{run}", deleteRun, "DELETE"},

	{"/triggers", listTriggers, "GET"},
	{"/triggers", addTrigger, "POST"},
	{"/triggers/{trigger}", getTrigger, "GET"},
	{"/triggers/{trigger}", updateTrigger, "PUT"},
	{"/triggers/{trigger}", deleteTrigger, "DELETE"},
	{"/triggers/{trigger}/jobs", listJobsForTrigger, "GET"},
}

type ctx struct {
	hub         *service.Hub
	executor    *service.Executor
	jobList     *service.JobList
	taskList    *service.TaskList
	triggerList *service.TriggerList
	runList     *service.RunList
}

func (t ctx) Hub() *service.Hub {
	return t.hub
}

func (t ctx) Executor() *service.Executor {
	return t.executor
}

func (t ctx) JobList() *service.JobList {
	return t.jobList
}

func (t ctx) TaskList() *service.TaskList {
	return t.taskList
}

func (t ctx) TriggerList() *service.TriggerList {
	return t.triggerList
}

func (t ctx) RunList() *service.RunList {
	return t.runList
}

type context interface {
	Hub() *service.Hub
	Executor() *service.Executor
	JobList() *service.JobList
	TaskList() *service.TaskList
	TriggerList() *service.TriggerList
	RunList() *service.RunList
}

type appHandler struct {
	*ctx
	handler func(context, http.ResponseWriter, *http.Request) (int, interface{})
}

func (t appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code, data := t.handler(t.ctx, w, r)
	marshal(data, w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	log.Println(r.URL, "-", r.Method, "-", code, r.RemoteAddr)
}

func main() {
	wd, _ := os.Getwd()
	log.Println("Working directory", wd)

	loadSettings("db/app.ini")
	log.Println("Writing logs to ", Settings.LogFile)
	if Settings.LogFile != "" {
		f, err := os.OpenFile(Settings.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("Can't open logfile! `%s` | Error: %#v\n", Settings.LogFile, err)
			os.Exit(1)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	jobList := service.NewJobList()
	taskList := service.NewTaskList()
	triggerList := service.NewTriggerList()
	runList := service.NewRunList(jobList)

	jobList.Load()
	taskList.Load()
	triggerList.Load()
	runList.Load()

	hub := service.NewHub(runList)
	go hub.HubLoop()

	executor := service.NewExecutor(jobList, taskList, runList)

	appContext := &ctx{hub, executor, jobList, taskList, triggerList, runList}

	r := mux.NewRouter()

	// non REST routes
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("web/")))
	r.HandleFunc("/", app).Methods("GET")
	r.HandleFunc("/favicon.ico", favicon).Methods("GET")
	r.Handle("/ws", appHandler{appContext, wsHandler}).Methods("GET")

	for _, detail := range routes {
		r.Handle(detail.route, appHandler{appContext, detail.handler}).Methods(detail.method)
	}

	host := Settings.ServeIP + ":" + strconv.Itoa(Settings.HTTPPort)

	log.Println("Running on " + host)
	http.ListenAndServe(host, r)
}
