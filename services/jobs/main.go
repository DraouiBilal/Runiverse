package main

import (
	"log"

	database "github.com/DraouiBilal/Runiverse-backend-lib/db"
	backend_service "github.com/DraouiBilal/Runiverse-backend-lib/service"
	job "github.com/DraouiBilal/Runiverse/jobs/service"
	"github.com/DraouiBilal/Runiverse/jobs/store/model"
)

func main() {

	db, err := database.ConnectDB("0.0.0.0", "postgres", "postgres", "runiverse-db", "disable", 5432)

	if err == nil {
		log.Println("DB connected")
	}else {
		log.Fatal("Can't connect to the DB")	
	}

	var jobService backend_service.CRUDInterface[model.Job] = &job.JobService{DB: db}

	_, db_err := database.CreateTable(db, "job", map[string][]string{
		"id":    {"VARCHAR(50)", "PRIMARY KEY"},
		"name":  {"VARCHAR(200)"},
	})

	if db_err != nil {
		log.Fatal("Can't create the table job", db_err)
	} else {
		log.Println("Table job created or already exists")
	}

	id := backend_service.GenerateID()

	job, create_err := jobService.Create(model.Job{
		Id:   id,
		Name: "My first job",
	})

	if create_err != nil {
		log.Fatal("Can't create a job:", create_err)
	} else {
		log.Println("Job created:", job)
	}
}

