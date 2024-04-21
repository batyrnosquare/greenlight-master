package data

import (
	"database/sql"
	"errors"
	"log"
)

type ModuleInfoModel struct {
	DB *sql.DB
}

func (m ModuleInfoModel) Insert(info *ModuleInfo) error {
	query := "INSERT INTO module_info(created_at, module_name, module_duration, exam_type) VALUES (now(), $1,$2,$3) RETURNING id, created_at, version"

	log.Println("inserted to db")

	return m.DB.QueryRow(
		query,
		&info.ModuleName,
		&info.ModuleDuration,
		&info.ExamType,
	).Scan(
		&info.ID,
		&info.CreatedAt,
		&info.Version,
	)
}

func (m ModuleInfoModel) Get(id int64) (*ModuleInfo, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := "SELECT module_name, module_duration, exam_type, version FROM module_info WHERE id = $1"

	var info ModuleInfo

	err := m.DB.QueryRow(query, id).Scan(
		&info.ModuleName,
		&info.ModuleDuration,
		&info.ExamType,
		&info.Version,
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

func (m ModuleInfoModel) GetLatestFifty() []*ModuleInfo {
	query := "SELECT * FROM module_info ORDER BY id DESC LIMIT 50"

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil
	}

	var infos []*ModuleInfo
	for rows.Next() {
		info := &ModuleInfo{}
		err = rows.Scan(
			&info.ID,
			&info.CreatedAt,
			&info.UpdatedAt,
			&info.ModuleName,
			&info.ModuleDuration,
			&info.ExamType,
			&info.Version,
		)
		if err != nil {
			return nil
		}
		infos = append(infos, info)
	}

	if err = rows.Err(); err != nil {
		return nil
	}
	return infos
}

func (m ModuleInfoModel) Update(info *ModuleInfo) error {
	query := "UPDATE module_info SET updated_at = now(), module_name = $1, module_duration = $2, exam_type = $3, version = version + 1 WHERE id = $4 RETURNING version"

	args := []interface{}{
		info.ModuleName,
		info.ModuleDuration,
		info.ExamType,
		info.ID,
	}
	return m.DB.QueryRow(query, args...).Scan(&info.Version)

}

func (m ModuleInfoModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := "DELETE FROM module_info WHERE id = $1"

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
