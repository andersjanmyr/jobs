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
	assert.Equal(t, "One", job.Name)
	assert.Equal(t, "one", job.Slug)
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
	j, _ := jobRepo.Add(job)
	assert.NotZero(t, j.ID)
	assert.NotZero(t, j.CreatedAt)
	assert.NotZero(t, j.UpdatedAt)
	jobs, _ := jobRepo.Find()
	assert.Equal(t, 3, len(jobs))
}

func TestUpdate(t *testing.T) {
	jobRepo := NewMemJobRepo([]*Job{NewJob("One"), NewJob("Two")})
	job, _ := jobRepo.FindOne("one")
	job.Name = "Uno"
	updatedAt := job.UpdatedAt
	updatedJob, _ := jobRepo.Update(job)
	assert.Equal(t, "Uno", updatedJob.Name)
	assert.NotEqual(t, updatedAt, updatedJob.UpdatedAt)
}

func TestUpdateFail(t *testing.T) {
	jobRepo := NewMemJobRepo([]*Job{NewJob("One"), NewJob("Two")})
	job := NewJob("Missing")
	job.Name = "Uno"
	j, err := jobRepo.Update(job)
	assert.Nil(t, j)
	assert.EqualError(t, err, "Cannot find job with slug missing")
}
