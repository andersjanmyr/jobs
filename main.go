package main

import (
	"encoding/json"
	"fmt"
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

	log.SetOutput(os.Stdout)
	var router = mux.NewRouter().StrictSlash(true)
	loggedRouter := handlers.LoggingHandler(os.Stdout, slowMiddleware(router))
	jobRepo := NewMemJobRepo([]*Job{NewJob("One"), NewJob("Two")})
	setupRouter(router.PathPrefix("/jobs"), NewJobController(jobRepo))

	go func() {
		log.Print("Profile server started on port 6060")
		log.Fatal(http.ListenAndServe("127.0.0.1:6060", nil))
	}()
	log.Print("Server started on port ", strconv.Itoa(port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), loggedRouter))
}

func slowMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 10000; i++ {
			fmt.Fprint(os.Stderr, i)
		}
		next.ServeHTTP(w, r)
	})
}

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

type RestController interface {
	Index(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Show(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Destroy(w http.ResponseWriter, r *http.Request)
	New(w http.ResponseWriter, r *http.Request)
	Edit(w http.ResponseWriter, r *http.Request)
}

type JobController struct {
	repo JobRepo
}

func NewJobController(repo JobRepo) *JobController {
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
	job, err := parseJob(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJson(w, job)
}

func parseJob(reader io.ReadCloser) (*Job, error) {
	if reader == nil {
		return nil, fmt.Errorf("No body to parse")
	}
	decoder := json.NewDecoder(reader)
	defer reader.Close() // errcheck-ignore
	var job Job
	if err := decoder.Decode(&job); err != nil {
		return nil, err
	}
	if job.Slug == "" {
		job.Slug = slug(job.Name)
	}
	return &job, nil
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
	job, err := parseJob(r.Body)
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
