package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/andersjanmyr/jobs/models"
	"github.com/gorilla/mux"
)

type JobController struct {
	repo models.JobRepo
}

func NewJobController(repo models.JobRepo) *JobController {
	jc := JobController{
		repo: repo,
	}
	return &jc
}

func (c *JobController) Index(w http.ResponseWriter, r *http.Request) {
	jobs, err := c.repo.Find()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJson(w, jobs)
}

func writeJson(w http.ResponseWriter, data interface{}) {
	json, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(json)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *JobController) Create(w http.ResponseWriter, r *http.Request) {
	job, err := models.ParseJob(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJson(w, job)
}

func (c *JobController) Show(w http.ResponseWriter, r *http.Request) {
	slug := getSlug(r)
	if slug == "" {
		http.NotFound(w, r)
		return
	}
	j, _ := c.repo.FindOne(slug)
	if j == nil {
		http.NotFound(w, r)
		return
	}
	writeJson(w, j)
}

func getSlug(r *http.Request) string {
	vars := mux.Vars(r)
	slug := vars["slug"]
	return slug
}

func (c *JobController) Update(w http.ResponseWriter, r *http.Request) {
	slug := getSlug(r)
	if slug == "" {
		http.NotFound(w, r)
		return
	}
	job, err := models.ParseJob(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	job.Slug = slug
	j, err := c.repo.UpAdd(job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJson(w, j)
}

func (c *JobController) Destroy(w http.ResponseWriter, r *http.Request) {
	slug := getSlug(r)
	if slug == "" {
		http.NotFound(w, r)
		return
	}
	j, err := c.repo.Delete(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	writeJson(w, j)
}
func (c *JobController) New(w http.ResponseWriter, r *http.Request)  {}
func (c *JobController) Edit(w http.ResponseWriter, r *http.Request) {}
