package mysql

import (
	"fmt"
	"time"
	"github.com/mihai22125/goPool/pkg/models"
	"database/sql"
)

// Define a PoolModel type which wraps a sql.DB connection pool
type SessionModel struct {
	DB *sql.DB
}

func (model *SessionModel) Insert(poolID, machineID int) (int, error) {
	stmt := `INSERT INTO session (pool_id, machine_id) VALUES (?, ?)`
	result, err := model.DB.Exec(stmt, poolID, machineID)
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

func (model *SessionModel) Get(id int) (*models.Session, error) {
	stmt := `SELECT session_id, pool_id, machine_id FROM session WHERE session_id = ?`
	session := &models.Session{}

	err := model.DB.QueryRow(stmt, id).Scan(&session.ID, &session.PoolID, &session.MachineID)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return session, nil
}

func (model *SessionModel) GetAll() ([]*models.Session, error) {
	stmt := `SELECT session_id, pool_id, machine_id FROM session`

	rows, err := model.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	
	sessions := []*models.Session{}

	for rows.Next() {
		session := &models.Session{}
	
		err = rows.Scan(&session.ID, &session.PoolID, &session.MachineID)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}


func (model *SessionModel) GetCurrentForMachine(machineID int) (*models.Session, error) {
	stmt := `SELECT session_id, session.pool_id, machine_id
			 FROM session
			 INNER JOIN pools on pools.pool_id = session.pool_id
			 INNER JOIN pool_config on pool_config.pool_id = pools.pool_id
			 WHERE machine_id = ? AND pool_config.start_date < ? AND pool_config.end_date > ?`

	session := &models.Session{}

	fmt.Println(time.Now())
	err := model.DB.QueryRow(stmt, machineID, time.Now().Format("2006-01-02 15:04:05") , time.Now().Format("2006-01-02 15:04:05") ).Scan(&session.ID, &session.PoolID, &session.MachineID)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	
	return session, nil
}
