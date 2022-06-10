package mysql

import (
	"database/sql"
)

// Define a PoolModel type which wraps a sql.DB connection pool
type VoteModel struct {
	DB *sql.DB
}

// Insert will insert a new pool int the database
func (model *VoteModel) Insert(poolID, optionID, machineID int, from [20]byte) (int, error) {
	stmt := `INSERT INTO vote (pool_id, option_id, machine_id, phone) VALUES (?, ?, ?, ?)`
	result, err := model.DB.Exec(stmt, poolID, optionID, machineID, from)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0,
			err
	}

	return int(id), nil
}

func (model *VoteModel) CountByOptionID(poolID, optionID int) (int, error) {
	rows, err := model.DB.Query("SELECT COUNT(*) FROM vote WHERE pool_id = ? AND option_id = ?", poolID, optionID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}

	return count, nil
}

func (model *VoteModel) CountByOptionIDDistinct(poolID, optionID int) (int, error) {
	rows, err := model.DB.Query("SELECT COUNT(DISTINCT phone) FROM vote WHERE pool_id = ? AND option_id = ?", poolID, optionID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}

	return count, nil
}
