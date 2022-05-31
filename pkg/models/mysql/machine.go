package mysql

import (
	"fmt"
	"time"
	"github.com/mihai22125/goPool/pkg/models"
	"database/sql"
)

// Define a PoolModel type which wraps a sql.DB connection pool
type MachineModel struct {
	DB *sql.DB
}

// Insert will insert a new pool int the database
func (model *MachineModel) Insert(machine models.Machine) (int, error) {
	stmt := `INSERT INTO machine (phone_number, ip_address) VALUES (?, ?)`
	result, err := model.DB.Exec(stmt, machine.PhoneNumber, machine.IPAdrres)
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

func (model *MachineModel) Get(id int) (*models.Machine, error) {
	stmt := `SELECT machine_id, phone_number, ip_address FROM machine WHERE machine_id = ?`
	machine := &models.Machine{}

	err := model.DB.QueryRow(stmt, id).Scan(&machine.ID, &machine.PhoneNumber, &machine.IPAdrres)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return machine, nil
}

// Update will update a specific pool based on its id
func (model *MachineModel) Update(machine *models.Machine) (int, error) {
	stmt := `UPDATE machine SET phone_number = ?, ip_address = ? WHERE machine_id = ?`

	_, err := model.DB.Exec(stmt, machine.PhoneNumber, machine.IPAdrres, machine.ID)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (model *MachineModel) Delete(id int) (int, error) {
	stmt := `DELETE FROM machine WHERE machine_id = ?`

	_, err := model.DB.Exec(stmt, id)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (model *MachineModel) GetAll() ([]*models.Machine, error) {
	stmt := `SELECT machine_id, phone_number, ip_address FROM machine`

	rows, err := model.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	
	machines := []*models.Machine{}

	for rows.Next() {
		machine := &models.Machine{}
	
		err = rows.Scan(&machine.ID, &machine.PhoneNumber, &machine.IPAdrres)
		if err != nil {
			return nil, err
		}

		machines = append(machines, machine)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return machines, nil
}

func (model *MachineModel) GetNextAvailable(start, end time.Time) (int, error) {
	stmt := `SELECT 
			 	m.machine_id, 
			    SUM(? <= pc.end_date AND ? >= pc.start_date) AS overlap_count
			 FROM machine as m
			 	INNER JOIN session as s ON s.machine_id = m.machine_id
			 	INNER JOIN pools as p ON p.pool_id = s.pool_id
				INNER JOIN pool_config as pc on pc.pool_id = p.pool_id
			 GROUP BY m.machine_id`

	rows, err := model.DB.Query(stmt, start, end)
	if err != nil {
		return 0, err
	}

	defer rows.Close()
	
	for rows.Next() {
		var machineID, count int
		fmt.Printf("%v  -- %v\n", machineID, count)
	
		err = rows.Scan(&machineID, &count)
		if err != nil {
			return 0, err
		}

		if count == 0 {
			return machineID, nil
		}
	}

// 	stmt := `SELECT 
// 	m.machine_id, 
//    pc.start_date,
//    pc.end_date,
//    (? <= pc.end_date) as now_before_end,
//    (? >= pc.start_date) as end_after_start
// FROM machine as m
// 	INNER JOIN session as s ON s.machine_id = m.machine_id
// 	INNER JOIN pools as p ON p.pool_id = s.pool_id
//    INNER JOIN pool_config as pc on pc.pool_id = p.pool_id
// `

// 	rows, err := model.DB.Query(stmt, start, end)
// 	if err != nil {
// 		return 0, err
// 	}

// 	defer rows.Close() 

// 	for rows.Next() {
// 		var machineID int
// 		var startDate, endDate time.Time
// 		var start1_before_end, end1_after_start bool
	
// 		_= rows.Scan(&machineID, &startDate, &endDate, &start1_before_end, &end1_after_start)
// 		fmt.Printf("%d: %v -- %v -- %v -- %v\n", machineID, startDate, endDate, start1_before_end, end1_after_start)
// 	}

	if err = rows.Err(); err != nil {
		return 0, err
	}

	return 0, nil
}
