package main

import (
	"fmt"
	"strings"
)

type Config struct {
}

type Job struct {
	Name   string
	Slug   string
	Config Config
}

func NewJob(name string) *Job {
	return &Job{name, slug(name), Config{}}
}

func slug(name string) string {
	return strings.ToLower(name)
}

func (j *Job) update(job *Job) {
	if job.Name != "" {
		j.Name = job.Name
	}
	j.Config = job.Config
}

type JobRepo struct {
	jobs []*Job
}

func NewJobRepo(jobs []*Job) *JobRepo {
	return &JobRepo{
		jobs: jobs,
	}
}

func (r *JobRepo) Find() []*Job {
	return r.jobs
}

func (r *JobRepo) FindOne(slug string) (*Job, int) {
	for i, j := range r.jobs {
		if j.Slug == slug {
			return j, i
		}
	}
	return nil, -1
}

func (r *JobRepo) Add(job *Job) *Job {
	r.jobs = append(r.jobs, job)
	return job
}

func (r *JobRepo) Update(job *Job) (*Job, error) {
	j, _ := r.FindOne(job.Slug)
	if j == nil {
		return nil, fmt.Errorf("Cannot find job with slug %s", job.Slug)
	}
	j.update(job)
	return job, nil
}

func (r *JobRepo) UpAdd(job *Job) *Job {
	j, _ := r.FindOne(job.Slug)
	if j == nil {
		return r.Add(job)
	} else {
		j.update(job)
		return j
	}
}

func (r *JobRepo) Delete(slug string) (*Job, error) {
	job, i := r.FindOne(slug)
	if job == nil {
		return nil, fmt.Errorf("Cannot find job with slug %s", job.Slug)
	}
	r.jobs = append(r.jobs[:i], r.jobs[i+1:]...)
	return job, nil
}
