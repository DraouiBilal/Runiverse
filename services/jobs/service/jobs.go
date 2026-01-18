package service

import (
	"database/sql"
	"fmt"
	database "github.com/DraouiBilal/Runiverse-backend-lib/db"
	"github.com/DraouiBilal/Runiverse/jobs/store/model"
)

type JobService struct{
	DB *sql.DB
}

func (jobService *JobService) GetAll() ([]model.Job, error) {
	job, err := database.QueryTable[model.Job](jobService.DB, "job")
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (jobService *JobService) GetById(id string) (model.Job, error) {
	row, err := database.Query(jobService.DB, "SELECT id, name FROM job WHERE id=$1", id)

	if err != nil {
		return model.Job{}, err
	}

	job := model.Job{}

	for row.Next() {
		if scan_err := row.Scan(&job.Id, &job.Name); scan_err != nil {
			return model.Job{}, scan_err
		}
	}

	return job, nil
}

func (jobService *JobService) Create(job model.Job) (model.Job, error) {
	_, err := database.Mutate(jobService.DB, fmt.Sprintf("INSERT INTO job(id, name) VALUES('job_%s', '%s')", job.Id , job.Name))
	if err != nil {
		return model.Job{}, err
	}
	return job, nil
}

func (jobService *JobService) Update(job model.Job) (model.Job, error) {
	_, err := database.Mutate(jobService.DB, fmt.Sprintf("UPDATE job SET name=%s WHERE id=%s", job.Name, job.Id))
	if err != nil {
		return model.Job{}, err
	}
	return job, nil
}

func (jobService *JobService) Delete(id string) error {
	_, err := database.Mutate(jobService.DB, fmt.Sprintf("DELETE job WHERE id=%s", id))
	if err != nil {
		return err
	}
	return nil
}
