package storage

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/calc"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/config"
)

type Scheduler struct {
	db     *sql.DB
	Config config.Config
}

var errNoTask = "no task with id = %s"

const limit = 10

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func New(db *sql.DB, config config.Config) Scheduler {
	return Scheduler{
		db:     db,
		Config: config,
	}
}

func (s *Scheduler) createDatabase(dbFile string) error {
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

	s.db = db
	return nil
}

func (s *Scheduler) GetTask(taskID string) (Task, error) {
	var task Task

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := s.db.QueryRow(query, taskID)

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return Task{}, fmt.Errorf(errNoTask, taskID)
		}
		return Task{}, err
	}

	return task, nil
}

func (s *Scheduler) GetTasks(search string) ([]Task, error) {
	var query string

	if search == "" {
		query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit`
	}
	if search != "" {
		searchDate, err := time.Parse("02.01.2006", search)
		if err == nil {
			search = searchDate.Format(calc.TimeFormat)
			query = `SELECT id, date, title, comment, repeat FROM scheduler
			 WHERE date LIKE :likePattern ORDER BY date LIMIT :limit`
		} else {
			query = `SELECT id, date, title, comment, repeat FROM scheduler 
			WHERE title LIKE :likePattern OR comment LIKE :likePattern ORDER BY date LIMIT :limit`
		}
	}

	rows, err := s.db.Query(query,
		sql.Named("limit", limit),
		sql.Named("likePattern", fmt.Sprintf("%%%s%%", search)),
	)
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

func (s *Scheduler) Update(task Task) error {
	query := `
	UPDATE scheduler 
	SET date = ?, title = ?, comment = ?, repeat = ?, id = ?
	WHERE id = ?;`
	rows, err := s.db.Exec(query, task.Date, &task.Title, task.Comment, task.Repeat, task.ID, task.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return fmt.Errorf(errNoTask, task.ID)

	}

	return nil
}

func (s *Scheduler) DeleteTask(taskID string) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	rows, err := s.db.Exec(query, taskID)
	if err != nil {
		return err
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return fmt.Errorf(errNoTask, taskID)

	}

	return nil
}

func (s *Scheduler) Add(task Task) (string, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"

	result, err := s.db.Exec(query, task.Date, &task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", fmt.Errorf("error executing query: %w", err)
	}

	resInt, err := (result.LastInsertId())
	if err != nil {
		return "", fmt.Errorf("error getting last id: %w", err)
	}
	res := strconv.FormatInt(resInt, 10)

	return res, nil
}

func (s *Scheduler) NewConnection() error {
	dbFile := s.Config.TODO_DBFILE
	fmt.Println(s.Config.TODO_DBFILE)

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return s.createDatabase(dbFile)
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("error opening db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging db: %w", err)
	}

	s.db = db
	return nil
}

func (t *Task) Validate() error {
	if t.Date == "" {
		t.Date = time.Now().Format(calc.TimeFormat)
	}

	if _, err := time.Parse(calc.TimeFormat, t.Date); err != nil {
		return fmt.Errorf("field `date` must be in format YYYYMMDD, but provided %w", err)
	}

	if t.Title == "" {
		return fmt.Errorf("field `title` cannot be empty")
	}

	nextDate, err := calc.NextDate(time.Now(), t.Date, t.Repeat)
	if err != nil {
		return fmt.Errorf("couldn't resolve next date: %w", err)
	}

	if t.Date < time.Now().Format(calc.TimeFormat) {
		if t.Repeat == "" {
			now := time.Now().Format(calc.TimeFormat)
			t.Date = now
		} else {
			t.Date = nextDate
		}
	}

	return nil
}
