package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PgJobRepo struct {
	db *gorm.DB
}

func NewPgJobRepo(db *gorm.DB) *PgJobRepo {
	return &PgJobRepo{
		db: db,
	}
}

func (r *PgJobRepo) Find() ([]*Job, error) {
	var jobs []*Job
	if err := r.db.Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *PgJobRepo) FindOne(slug string) (*Job, error) {
	job := Job{}
	if err := r.db.Where(&Job{Slug: slug}).First(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *PgJobRepo) Add(job *Job) (*Job, error) {
	if err := r.db.Create(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

func (r *PgJobRepo) Update(job *Job) (*Job, error) {
	newJob := *job
	if err := r.db.Save(&newJob).Error; err != nil {
		return nil, err
	}
	return &newJob, nil
}

func (r *PgJobRepo) UpAdd(job *Job) (*Job, error) {
	existingJob, err := r.FindOne(job.Slug)
	if err != nil {
		return r.Add(job)
	}
	return r.Update(existingJob)

}

func (r *PgJobRepo) Delete(slug string) (*Job, error) {
	job, err := r.FindOne(slug)
	if err != nil {
		return nil, err
	}
	if err := r.db.Delete(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}
