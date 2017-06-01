package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	jobRepo := NewMemJobRepo([]*Job{NewJob("One"), NewJob("Two")})
	jobs, _ := jobRepo.Find()
	assert.Equal(t, 2, len(jobs))
}

func TestFindOne(t *testing.T) {
	jobRepo := NewMemJobRepo([]*Job{NewJob("One"), NewJob("Two")})
	job, err := jobRepo.FindOne("one")
	assert.Nil(t, err)
	assert.Equal(t, NewJob("One"), job)
}

func TestFindOneMissing(t *testing.T) {
	jobRepo := NewMemJobRepo([]*Job{NewJob("One"), NewJob("Two")})
	job, err := jobRepo.FindOne("missing")
	assert.Nil(t, job)
	assert.EqualError(t, err, "No job found with slug: missing")
}

func TestAdd(t *testing.T) {
	jobRepo := NewMemJobRepo([]*Job{NewJob("One"), NewJob("Two")})
	job := NewJob("One")
	_, _ = jobRepo.Add(job)
	jobs, _ := jobRepo.Find()
	assert.Equal(t, 3, len(jobs))
}

func TestUpdate(t *testing.T) {
	jobRepo := NewMemJobRepo([]*Job{NewJob("One"), NewJob("Two")})
	job := NewJob("One")
	job.Name = "Uno"
	updatedJob, _ := jobRepo.Update(job)
	assert.Equal(t, "Uno", updatedJob.Name)
}

func TestUpdateFail(t *testing.T) {
	jobRepo := NewMemJobRepo([]*Job{NewJob("One"), NewJob("Two")})
	job := NewJob("Missing")
	job.Name = "Uno"
	j, err := jobRepo.Update(job)
	assert.Nil(t, j)
	assert.EqualError(t, err, "Cannot find job with slug missing")
}
