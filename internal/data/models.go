package data

import (
	"database/sql"
	"errors"
	"time"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record (row, entry) not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Create a Models struct which wraps the MovieModel
// kind of enveloping
type Models struct {
	Movies      MovieModel
	ModuleInfos ModuleInfoModel
	DepInfos    DepartmentInfoModel
	UserInfos   UserInfoModel
	Tokens      TokenModel
}

// method which returns a Models struct containing the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Movies:      MovieModel{DB: db},
		DepInfos:    DepartmentInfoModel{DB: db},
		ModuleInfos: ModuleInfoModel{DB: db},
		UserInfos:   UserInfoModel{DB: db},
		Tokens:      TokenModel{DB: db},
	}
}

type ModuleInfo struct {
	ID             int           `json:"id,omitempty"`
	CreatedAt      time.Time     `json:"-"`
	UpdatedAt      time.Time     `json:"updatedAt,omitempty"`
	ModuleName     string        `json:"moduleName"`
	ModuleDuration time.Duration `json:"moduleDuration"`
	ExamType       string        `json:"examType"`
	Version        string        `json:"version"`
}

type DepartmentInfo struct {
	ID                 int    `json:"id"`
	DepartmentName     string `json:"departmentName"`
	DepartmentDirector string `json:"departmentDirector"`
	StaffQuantity      int    `json:"staffQuantity"`
	ModuleID           int    `json:"moduleId"`
}

type UserInfo struct {
	ID           int       `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"passwordHash"`
	Role         string    `json:"role"`
	Activated    bool      `json:"activated"`
	Version      int       `json:"version"`
}

const (
	Admin      = "admin"
	Registered = "registered"
)
