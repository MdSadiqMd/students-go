package storage

import (
	"github.com/MdSadiqMd/students-go/internal/types"
)

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id string) (types.Student, error)
}
