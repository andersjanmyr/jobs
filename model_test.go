package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	jobRepo := NewJobRepo([]*Job{NewJob("One"), NewJob("Two")})
	assert.Equal(t, 2, len(jobRepo.Find()))
}

func TestFindOne(t *testing.T) {
	jobRepo := NewJobRepo([]*Job{NewJob("One"), NewJob("Two")})
	job, i := jobRepo.FindOne("one")
	assert.Equal(t, NewJob("One"), job)
	assert.Equal(t, i, 0)
}

func TestFindOneMissing(t *testing.T) {
	jobRepo := NewJobRepo([]*Job{NewJob("One"), NewJob("Two")})
	job, i := jobRepo.FindOne("missing")
	assert.Nil(t, job)
	assert.Equal(t, i, -1)
}

func TestAdd(t *testing.T) {
	jobRepo := NewJobRepo([]*Job{NewJob("One"), NewJob("Two")})
	job := NewJob("One")
	jobRepo.Add(job)
	assert.Equal(t, 3, len(jobRepo.Find()))
}

func TestUpdate(t *testing.T) {
	jobRepo := NewJobRepo([]*Job{NewJob("One"), NewJob("Two")})
	job := NewJob("One")
	job.Name = "Uno"
	updatedJob, _ := jobRepo.Update(job)
	assert.Equal(t, "Uno", updatedJob.Name)
}

func TestUpdateFail(t *testing.T) {
	jobRepo := NewJobRepo([]*Job{NewJob("One"), NewJob("Two")})
	job := NewJob("Missing")
	job.Name = "Uno"
	j, err := jobRepo.Update(job)
	assert.Nil(t, j)
	assert.EqualError(t, err, "Cannot find job with slug missing")
}
