package data

import (
	"database/sql"
	"errors"
	"log"
)

type DepartmentInfoModel struct {
	DB *sql.DB
}

func (m DepartmentInfoModel) Insert(info *DepartmentInfo) error {
	query := "INSERT INTO department_info(department_name, department_director, staff_quantity, module_id) VALUES ($1,$2,$3, $4) RETURNING id "

	log.Println("inserted to db")

	return m.DB.QueryRow(
		query,
		&info.DepartmentName,
		&info.DepartmentDirector,
		&info.StaffQuantity,
		&info.ModuleID,
	).Scan(
		&info.ID,
	)
}

func (m DepartmentInfoModel) Get(id int64) (*DepartmentInfo, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := "SELECT * FROM department_info WHERE id = $1"

	var info DepartmentInfo

	err := m.DB.QueryRow(query, id).Scan(
		&info.ID,
		&info.DepartmentName,
		&info.DepartmentDirector,
		&info.StaffQuantity,
		&info.ModuleID,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err

		}
	}
	return &info, nil
}
