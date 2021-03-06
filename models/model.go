package models

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Job struct {
	gorm.Model
	Name string
	Slug string
}

var jobIndex uint = 0

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
	if job.Slug != "" {
		j.Slug = job.Slug
	}
	job.UpdatedAt = time.Now()
}

func ParseJob(reader io.ReadCloser) (*Job, error) {
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

func ParseJobs(reader io.ReadCloser) ([]Job, error) {
	if reader == nil {
		return nil, fmt.Errorf("No body to parse")
	}
	decoder := json.NewDecoder(reader)
	defer reader.Close()
	var jobs []Job
	if err := decoder.Decode(&jobs); err != nil {
		return nil, err
	}
	return jobs, nil
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
	r := &MemJobRepo{
		jobs: []*Job{},
	}
	for _, j := range jobs {
		_, _ = r.Add(j)
	}
	return r
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
	jobIndex++
	now := time.Now()
	job.CreatedAt = now
	job.UpdatedAt = now
	job.ID = jobIndex
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
