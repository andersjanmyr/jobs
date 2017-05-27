package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "net/http/pprof"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := 5555

	r := mux.NewRouter()
	r.HandleFunc("/", pipeline)
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	jobs := r.PathPrefix("/jobs").Subrouter()
	jobsController := newJobsController()
	jobs.HandleFunc("/", jobsController.Index).Methods("GET")
	jobs.HandleFunc("/", jobsController.Create).Methods("POST")
	jobs.HandleFunc("/{name}", jobsController.Show).Methods("GET")
	jobs.HandleFunc("/{name}", jobsController.Update).Methods("PUT")
	jobs.HandleFunc("/{name}", jobsController.Destroy).Methods("DELETE")
	jobs.HandleFunc("/{name}/new", jobsController.New).Methods("GET")
	jobs.HandleFunc("/{name}/edit", jobsController.Edit).Methods("GET")

	log.Print("Server started on port; ", strconv.Itoa(port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), loggedRouter))
}

func pipeline(w http.ResponseWriter, r *http.Request) {
	log.Print("Pipeline")
	w.Write([]byte("pipeline"))
}

type Job struct {
	Name string
}

type RestController interface {
	Index(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Show(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Destroy(w http.ResponseWriter, r *http.Request)
	New(w http.ResponseWriter, r *http.Request)
	Edit(w http.ResponseWriter, r *http.Request)
}

type JobsController struct {
	Jobs []*Job
}

func newJobsController() JobsController {
	jc := JobsController{
		Jobs: []*Job{
			&Job{"One"},
			&Job{"Two"},
		},
	}
	return jc
}

func (c *JobsController) Index(w http.ResponseWriter, r *http.Request) {
	json, err := json.MarshalIndent(c.Jobs, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(json)
}

func (c *JobsController) Create(w http.ResponseWriter, r *http.Request) {
	job, err := parseJob(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Jobs = append(c.Jobs, job)
	json, err := json.MarshalIndent(job, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(json)
}

func parseJob(reader io.ReadCloser) (*Job, error) {
	decoder := json.NewDecoder(reader)
	defer reader.Close()
	var job Job
	if err := decoder.Decode(&job); err != nil {
		return nil, err
	}
	return &job, nil
}

func (c *JobsController) Show(w http.ResponseWriter, r *http.Request) {
	name := getName(r)
	if name == "" {
		http.NotFound(w, r)
		return
	}
	j := c.findJob(name)
	if j == nil {
		http.NotFound(w, r)
	}
	json, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func getName(r *http.Request) string {
	vars := mux.Vars(r)
	name := vars["name"]
	log.Println("name", name)
	return name
}

func (c *JobsController) Update(w http.ResponseWriter, r *http.Request) {
	name := getName(r)
	if name == "" {
		http.NotFound(w, r)
		return
	}
	job, err := parseJob(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	j := c.findJob(name)
	if j == nil {
		c.Jobs = append(c.Jobs, job)
	} else {
		// ignore for now
	}
	json, err := json.MarshalIndent(job, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func (c *JobsController) findJob(name string) *Job {
	for _, j := range c.Jobs {
		if j.Name == name {
			return j
		}
	}
	return nil
}

func (c *JobsController) Destroy(w http.ResponseWriter, r *http.Request) {}
func (c *JobsController) New(w http.ResponseWriter, r *http.Request)     {}
func (c *JobsController) Edit(w http.ResponseWriter, r *http.Request)    {}

func JobsHandler(w http.ResponseWriter, r *http.Request) {
}
