package store

import (
	"database/sql"
	"fmt"
	database "github.com/DraouiBilal/Runiverse-backend-lib/db"
	"github.com/DraouiBilal/Runiverse/history/model"
)

func Get(db *sql.DB) ([]model.History, error) {
	history, err := database.QueryTable[model.History](db, "history")
	if err != nil {
		return nil, err
	}
	return history, nil
}

func GetById(db *sql.DB, id string) (model.History, error) {
	row, err := database.Query(db, "SELECT id, status, logs FROM history WHERE id=$1", id)

	if err != nil {
		return model.History{}, err
	}

	history := model.History{}

	for row.Next() {
		if scan_err := row.Scan(&history.Id, &history.Status, &history.Logs); scan_err != nil {
			return model.History{}, scan_err
		}
	}

	return history, nil
}

func Create(db *sql.DB, history model.History) (model.History, error) {
	_, err := database.Mutate(db, fmt.Sprintf("INSERT INTO history(status, logs) VALUES('%s', '%s')", history.Status, history.Logs))
	if err != nil {
		return model.History{}, err
	}
	return history, nil
}

func Update(db *sql.DB, history model.History) (model.History, error) {
	_, err := database.Mutate(db, fmt.Sprintf("UPDATE history SET status=%s, logs=%s WHERE id=%s", history.Status, history.Logs, history.Id))
	if err != nil {
		return model.History{}, err
	}
	return history, nil
}

func Delete(db *sql.DB, history model.History) error {
	_, err := database.Mutate(db, fmt.Sprintf("DELETE history WHERE id=%s", history.Id))
	if err != nil {
		return err
	}
	return nil
}
