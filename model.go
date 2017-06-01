package main

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

type Job struct {
	gorm.Model
	Name string
	Slug string
}

func NewJob(name string) *Job {
	return &Job{Name: name, Slug: slug(name)}
}

func slug(name string) string {
	return strings.ToLower(name)
}

func (j *Job) update(job *Job) {
	if job.Name != "" {
		j.Name = job.Name
	}
}

type JobRepo interface {
	Find() ([]*Job, error)
	FindOne(slug string) (*Job, error)
	Add(job *Job) (*Job, error)
	Update(job *Job) (*Job, error)
	UpAdd(job *Job) (*Job, error)
	Delete(slug string) (*Job, error)
}

type MemJobRepo struct {
	jobs []*Job
}

func NewMemJobRepo(jobs []*Job) *MemJobRepo {
	return &MemJobRepo{
		jobs: jobs,
	}
}

func (r *MemJobRepo) Find() ([]*Job, error) {
	return r.jobs, nil
}

func (r *MemJobRepo) FindOne(slug string) (*Job, error) {
	for _, j := range r.jobs {
		if j.Slug == slug {
			return j, nil
		}
	}
	return nil, fmt.Errorf("No job found with slug: %s", slug)
}

func (r *MemJobRepo) Add(job *Job) (*Job, error) {
	r.jobs = append(r.jobs, job)
	return job, nil
}

func (r *MemJobRepo) index(slug string) int {
	for i, j := range r.jobs {
		if j.Slug == slug {
			return i
		}
	}
	return -1
}
func (r *MemJobRepo) Update(job *Job) (*Job, error) {
	j, _ := r.FindOne(job.Slug)
	if j == nil {
		return nil, fmt.Errorf("Cannot find job with slug %s", job.Slug)
	}
	j.update(job)
	return job, nil
}

func (r *MemJobRepo) UpAdd(job *Job) (*Job, error) {
	j, _ := r.FindOne(job.Slug)
	if j == nil {
		return r.Add(job)
	} else {
		j.update(job)
		return j, nil
	}
}

func (r *MemJobRepo) Delete(slug string) (*Job, error) {
	i := r.index(slug)
	if i == -1 {
		return nil, fmt.Errorf("Cannot find job with slug %s", slug)
	}
	job := r.jobs[i]
	r.jobs = append(r.jobs[:i], r.jobs[i+1:]...)
	return job, nil
}
