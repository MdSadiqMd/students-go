package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/MdSadiqMd/students-go/internal/config"
	"github.com/MdSadiqMd/students-go/internal/types"
	_ "github.com/mattn/go-sqlite3" // we are using this package indirectly (not using function from it) so to avoid error we use _ blank identifier
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) { // constructor which returns instance of Sqlite
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS students (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			age INTEGER NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}

	return &Sqlite{Db: db}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, nil
		}
		return types.Student{}, fmt.Errorf("student not found: %w", err)
	}
	return student, nil
}

func (s *Sqlite) GetStudentList() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	var students []types.Student
	for rows.Next() {
		var student types.Student
		err = rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}

func (s *Sqlite) UpdateStudent(id int64, name string, email string, age int) error {
	var query string
	var values []interface{}

	if name != "" {
		query += "name = ?,"
		values = append(values, name)
	}
	if email != "" {
		query += "email = ?,"
		values = append(values, email)
	}
	if age != 0 {
		query += "age = ?,"
		values = append(values, age)
	}

	if query == "" {
		return nil
	}

	query = "UPDATE students SET " + query[:len(query)-1] + " WHERE id = ?"
	stmt, err := s.Db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	values = append(values, id)
	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}
	return nil
}
