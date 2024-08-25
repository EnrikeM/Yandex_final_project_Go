package storage

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (dbParams *DBParams) createDatabase(dbFile string) error {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("error opening db: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date   VARCHAR(32) NOT NULL DEFAULT "0001-01-01",
		title   VARCHAR(64) NOT NULL DEFAULT "",
		comment VARCHAR(128) NOT NULL DEFAULT "",
		repeat  VARCHAR(128) NOT NULL DEFAULT ""
	);
	CREATE INDEX scheduler_date ON scheduler (date);`)
	if err != nil {
		return fmt.Errorf("error creating db: %w", err)
	}

	dbParams.DB = db
	return nil
}

func GetTask(db *sql.DB, taskID string) (Task, error) { //вынести в utils?
	var task Task

	query := "SELECT * FROM scheduler WHERE id = ?"
	row := db.QueryRow(query, taskID)

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func GetTasks(db *sql.DB, search string) ([]Task, error) { //вынести в utils?
	var query string

	if search == "" {
		query = "SELECT * FROM scheduler ORDER BY date DESC LIMIT 10"
	}
	if search != "" {
		searchDate, err := time.Parse("02.01.2006", search)
		if err == nil {
			search = searchDate.Format("20060102")
			query = "SELECT * FROM scheduler WHERE date LIKE ? ORDER BY date DESC LIMIT 10"
		} else {
			query = "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT 10"
		}
	}

	rows, err := db.Query(query, fmt.Sprintf("%%%s%%", search), fmt.Sprintf("%%%s%%", search))
	if err != nil {
		return nil, fmt.Errorf("error exectuing query: %v", err)
	}

	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return tasks, nil
}

func RedactTask(db *sql.DB, task Task) error {
	query := `
	UPDATE scheduler 
	SET date = ?, title = ?, comment = ?, repeat = ?, id = ?
	WHERE id = ?;`
	_, err := db.Exec(query, task.Date, &task.Title, task.Comment, task.Repeat, task.ID, task.ID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTask(db *sql.DB, taskID string) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	_, err := db.Exec(query, taskID)
	if err != nil {
		return err
	}
	return nil
}

func GetLastId(task Task, db *sql.DB) (string, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"

	result, err := db.Exec(query, task.Date, &task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", fmt.Errorf("error executing query: %w", err)
	}

	resInt, err := (result.LastInsertId())
	if err != nil {
		return "", fmt.Errorf("error getting last id: %w", err)
	}
	res := strconv.Itoa(int(resInt))

	return res, nil
}
