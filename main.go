package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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
	Jobs []Job
}

func newJobsController() JobsController {
	jc := JobsController{
		Jobs: []Job{
			Job{"One"},
			Job{"Two"},
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
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var job Job
	err := decoder.Decode(&job)
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

func (c *JobsController) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	log.Println("name", name)
	fmt.Println("name", name)
	if name == "" {
		http.NotFound(w, r)
		return
	}
	for _, j := range c.Jobs {
		if j.Name == name {
			json, err := json.MarshalIndent(j, "", "  ")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(json)
			return
		}
		http.NotFound(w, r)
	}

}
func (c *JobsController) Update(w http.ResponseWriter, r *http.Request)  {}
func (c *JobsController) Destroy(w http.ResponseWriter, r *http.Request) {}
func (c *JobsController) New(w http.ResponseWriter, r *http.Request)     {}
func (c *JobsController) Edit(w http.ResponseWriter, r *http.Request)    {}

func JobsHandler(w http.ResponseWriter, r *http.Request) {
}
