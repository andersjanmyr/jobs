package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestJobsIndex(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter()
	controller := newJobsController()
	setupRouter(router.PathPrefix("/"), controller)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	expected := `[
		  {
		    "Name": "One",
		    "Slug": "one",
		    "Config": {}
		  },
		  {
		    "Name": "Two",
		    "Slug": "two",
		    "Config": {}
		  }
		]`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestJobsCreate(t *testing.T) {
	job := strings.NewReader(`{
		"Name": "Three"
	}`)

	req, err := http.NewRequest("POST", "/", job)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter()
	controller := newJobsController()
	setupRouter(router.PathPrefix("/"), controller)
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	expected := `{
		"Name": "Three",
		"Slug": "three",
		"Config": {}
	}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestJobsShow(t *testing.T) {
	req, err := http.NewRequest("GET", "/one", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter()
	controller := newJobsController()
	controller.Jobs = []*Job{newJob("Zero"), newJob("One"), newJob("Two")}
	setupRouter(router.PathPrefix("/"), controller)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	expected := `{
		"Name": "One",
		"Slug": "one",
		"Config": {}
	}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestJobsUpdate(t *testing.T) {
	job := strings.NewReader(`{
		"Name": "Uno"
	}`)

	req, err := http.NewRequest("PUT", "/one", job)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter()
	controller := newJobsController()
	controller.Jobs = []*Job{newJob("one")}
	setupRouter(router.PathPrefix("/"), controller)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	expected := `{
		"Name": "Uno",
		"Slug": "one",
		"Config": {}
	}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestJobsDelete(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/one", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter()
	controller := newJobsController()
	controller.Jobs = []*Job{newJob("Zero"), newJob("One"), newJob("Two")}
	setupRouter(router.PathPrefix("/"), controller)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	expected := `{
		"Name": "One",
		"Slug": "one",
		"Config": {}
	}`
	assert.JSONEq(t, expected, w.Body.String())
	assert.Equal(t, 2, len(controller.Jobs))
}
