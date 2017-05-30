package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func jsonToMap(w *httptest.ResponseRecorder) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &m)
	return m
}

func jsonToSlice(w *httptest.ResponseRecorder) []map[string]interface{} {
	var s []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &s)
	return s
}

func TestJobsIndex(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter()
	jobRepo := NewJobRepo([]*Job{NewJob("One"), NewJob("Two")})
	controller := NewJobController(jobRepo)
	setupRouter(router.PathPrefix("/"), controller)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	s := jsonToSlice(w)
	assert.Equal(t, 2, len(s))
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
	controller := NewJobController(NewJobRepo([]*Job{}))
	setupRouter(router.PathPrefix("/"), controller)
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	m := jsonToMap(w)
	assert.Equal(t, "Three", m["Name"])
	assert.Equal(t, "three", m["Slug"])
}

func TestJobsShow(t *testing.T) {
	req, err := http.NewRequest("GET", "/one", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter()
	jobRepo := NewJobRepo([]*Job{NewJob("Zero"), NewJob("One"), NewJob("Two")})
	controller := NewJobController(jobRepo)
	setupRouter(router.PathPrefix("/"), controller)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	m := jsonToMap(w)
	assert.Equal(t, "One", m["Name"])
	assert.Equal(t, "one", m["Slug"])
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
	jobRepo := NewJobRepo([]*Job{NewJob("One"), NewJob("Two")})
	controller := NewJobController(jobRepo)
	setupRouter(router.PathPrefix("/"), controller)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	m := jsonToMap(w)
	assert.Equal(t, "Uno", m["Name"])
}

func TestJobsDelete(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/one", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter()
	jobRepo := NewJobRepo([]*Job{NewJob("Zero"), NewJob("One"), NewJob("Two")})
	controller := NewJobController(jobRepo)
	setupRouter(router.PathPrefix("/"), controller)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	m := jsonToMap(w)
	assert.Equal(t, "One", m["Name"])
	assert.Equal(t, 2, len(jobRepo.Find()))
}
