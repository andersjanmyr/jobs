package main

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	db, err = gorm.Open("postgres", "host=localhost user=jobs dbname=jobs sslmode=disable password=jobs")
	if err != nil {
		panic(err)
	}
	defer db.Close() // errcheck-ignore

	db.AutoMigrate(&Job{})
	code := m.Run()

	os.Exit(code)
}

func TestPgJobsAddFind(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()
	jobRepo := NewPgJobRepo(tx)
	job := NewJob("Dingo")
	job, err := jobRepo.Add(job)
	assert.Nil(t, err)
	jobs, err := jobRepo.Find()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(jobs))
}

func TestPgJobsUpdate(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()
	jobRepo := NewPgJobRepo(tx)
	job := NewJob("Dingo")
	job, err := jobRepo.Add(job)
	job.Name = "Sloth"
	updatedJob, err := jobRepo.Update(job)
	assert.Nil(t, err)
	assert.True(t, updatedJob.UpdatedAt.After(job.UpdatedAt), "Expected '%v' to be after '%v'", updatedJob.UpdatedAt, job.UpdatedAt)
}

func TestPgJobsUpAdd(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()
	jobRepo := NewPgJobRepo(tx)
	job := NewJob("Dingo")
	job, err := jobRepo.UpAdd(job)
	assert.Nil(t, err)
	assert.NotNil(t, job)
}

func TestPgJobsFindOne(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()
	jobRepo := NewPgJobRepo(tx)
	_, _ = jobRepo.Add(NewJob("Dingo"))
	_, _ = jobRepo.Add(NewJob("Sloth"))
	job, err := jobRepo.FindOne("sloth")
	assert.Nil(t, err)
	assert.Equal(t, "Sloth", job.Name)
}

func TestPgJobsDelete(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()
	jobRepo := NewPgJobRepo(tx)
	_, _ = jobRepo.Add(NewJob("Dingo"))
	_, _ = jobRepo.Add(NewJob("Sloth"))
	job, err := jobRepo.Delete("sloth")
	assert.Nil(t, err)
	assert.Equal(t, "Sloth", job.Name)
	jobs, err := jobRepo.Find()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(jobs))
}
