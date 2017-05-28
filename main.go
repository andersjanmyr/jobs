package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "net/http/pprof"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func setupRouter(router *mux.Route, controller RestController) *mux.Router {
	var subRouter = router.Subrouter()
	subRouter.HandleFunc("/", controller.Index).Methods("GET")
	subRouter.HandleFunc("/", controller.Create).Methods("POST")
	subRouter.HandleFunc("/{slug}", controller.Show).Methods("GET")
	subRouter.HandleFunc("/{slug}", controller.Update).Methods("PUT")
	subRouter.HandleFunc("/{slug}", controller.Destroy).Methods("DELETE")
	subRouter.HandleFunc("/{slug}/new", controller.New).Methods("GET")
	subRouter.HandleFunc("/{slug}/edit", controller.Edit).Methods("GET")
	return subRouter
}

func main() {
	port := 5555

	var router = mux.NewRouter()
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	setupRouter(router.PathPrefix("/jobs"), newJobsController())

	log.Print("Server started on port; ", strconv.Itoa(port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), loggedRouter))
}

type Config struct {
}

type Job struct {
	Name   string
	Slug   string
	Config Config
}

func newJob(name string) *Job {
	return &Job{name, slug(name), Config{}}
}

func slug(name string) string {
	return strings.ToLower(name)
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

func newJobsController() *JobsController {
	jc := JobsController{
		Jobs: []*Job{
			newJob("One"),
			newJob("Two"),
		},
	}
	return &jc
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
	if reader == nil {
		return nil, fmt.Errorf("No body to parse")
	}
	decoder := json.NewDecoder(reader)
	defer reader.Close()
	var job Job
	if err := decoder.Decode(&job); err != nil {
		return nil, err
	}
	if job.Slug == "" {
		job.Slug = slug(job.Name)
	}
	return &job, nil
}

func (c *JobsController) Show(w http.ResponseWriter, r *http.Request) {
	slug := getSlug(r)
	if slug == "" {
		http.NotFound(w, r)
		return
	}
	j, _ := c.findJob(slug)
	if j == nil {
		http.NotFound(w, r)
		return
	}
	json, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func getSlug(r *http.Request) string {
	vars := mux.Vars(r)
	slug := vars["slug"]
	return slug
}

func (c *JobsController) Update(w http.ResponseWriter, r *http.Request) {
	slug := getSlug(r)
	if slug == "" {
		http.NotFound(w, r)
		return
	}
	job, err := parseJob(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	j, _ := c.findJob(slug)
	if j == nil {
		c.Jobs = append(c.Jobs, job)
	} else {
		if job.Name != "" {
			j.Name = job.Name
		}
		j.Config = job.Config
	}
	json, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func (c *JobsController) findJob(slug string) (*Job, int) {
	for i, j := range c.Jobs {
		if j.Slug == slug {
			return j, i
		}
	}
	return nil, -1
}

func (c *JobsController) Destroy(w http.ResponseWriter, r *http.Request) {
	slug := getSlug(r)
	if slug == "" {
		http.NotFound(w, r)
		return
	}
	j, i := c.findJob(slug)
	if j == nil {
		http.NotFound(w, r)
		return
	}
	c.Jobs = append(c.Jobs[:i], c.Jobs[i+1:]...)
	json, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}
func (c *JobsController) New(w http.ResponseWriter, r *http.Request)  {}
func (c *JobsController) Edit(w http.ResponseWriter, r *http.Request) {}
