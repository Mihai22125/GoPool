package mysql

import (
	"context"
	"github.com/mihai22125/goPool/pkg/models"
	"database/sql"
)

// Define a PoolModel type which wraps a sql.DB connection pool
type PoolModel struct {
	DB *sql.DB
}

// Insert will insert a new pool int the database
func (model *PoolModel) Insert(pool models.Pool) (int, error) {
	stmt := `INSERT INTO pools (user_id, name, nr_of_options) VALUES (?, ?, ?)`
	result, err := model.DB.Exec(stmt, pool.UserID, pool.Name, pool.NumberOfOptions)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, 
		err
	}

	pool.ID = int(id)
	pool.PoolConfig.PoolID = int(id)
	_, err = model.InsertConfig(pool)
	if err != nil {
		return 0, 
		err
	}

	return int(id), nil
}

// Get will return a specific pool based on its id
func (model *PoolModel) Get(id int) (*models.Pool, error) {
	stmt := `SELECT pool_id, user_id, name, nr_of_options FROM pools WHERE pool_id = ?`
	pool := &models.Pool{}

	err := model.DB.QueryRow(stmt, id).Scan(&pool.ID, &pool.UserID, &pool.Name, &pool.NumberOfOptions)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	pool.PoolOptions, err = model.GetOptions(*pool)
    if err != nil {
		return nil, err
	}

	pool.PoolConfig, err = model.GetConfig(*pool)
	if err != nil && err != models.ErrNoRecord {
		return nil, err
	}

	return pool, nil
}

// Update will update a specific pool based on its id
func (model *PoolModel) Update(pool *models.Pool) (int, error) {
	stmt := `UPDATE pools SET name = ?, nr_of_options = ? WHERE pool_id = ?`

	_, err := model.DB.Exec(stmt, pool.Name, pool.NumberOfOptions, pool.ID)
	if err != nil {
		return 0, err
	}

	_, err = model.GetConfig(*pool)
	if err == models.ErrNoRecord {
		_, err = model.InsertConfig(*pool)
		if err != nil {
			return 0, err
		}
	} else {
		_, err = model.UpdateConfig(*pool)
		if err != nil {
			return 0, err
		}
	}

	return 0, nil
}

// GetAll will return all created pools
func (model *PoolModel) GetAll(userID int) ([]*models.Pool, error) {
	stmt := `SELECT pool_id, user_id, name, nr_of_options FROM pools where user_id = ?`

	rows, err := model.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	
	pools := []*models.Pool{}

	for rows.Next() {
		pool := &models.Pool{}
	
		err = rows.Scan(&pool.ID, &pool.UserID, &pool.Name, &pool.NumberOfOptions)
		if err != nil {
			return nil, err
		}

		pool.PoolOptions, err = model.GetOptions(*pool)
		if err != nil {
			return nil, err
		}
	
		pool.PoolConfig, err = model.GetConfig(*pool)
		if err != nil && err != models.ErrNoRecord {
			return nil, err
		}

		pools = append(pools, pool)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pools, nil
}


func (model *PoolModel) InsertOption(option models.PoolOption) (int, error) {
	stmt := `INSERT INTO pool_option (pool_id, name, description) VALUES (?, ?, ?)`
	result, err := model.DB.Exec(stmt, option.PoolID, option.Option, option.Description)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	} 

	return int(id), nil
}

func (model *PoolModel) GetOptions(pool models.Pool) ([]models.PoolOption, error) {
	stmt := `SELECT option_id, pool_id, name, description FROM pool_option where pool_id = ?`

	rows, err := model.DB.Query(stmt, pool.ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	options := make([]models.PoolOption, pool.NumberOfOptions)

	if pool.NumberOfOptions == 0 {
		return options, nil
	}
	
	var i int
	for rows.Next() {
		err = rows.Scan(&options[i].ID, &options[i].PoolID, &options[i].Option, &options[i].Description)
		if err != nil {
			return nil, err
		}

		i++
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return options, nil
}

func (model *PoolModel) DeleteAllOptions(poolID int) (int, error) {
	stmt := `DELETE from pool_option where pool_id = ?`
	result, err := model.DB.Exec(stmt, poolID)
	if err != nil {
		return 0, err
	}

	id, err := result.RowsAffected()
	if err != nil {
		return 0, err
	} 

	return int(id), nil
}

// TODO: make a mysql transaction
func (model *PoolModel) UpdateOptions(pool models.Pool) (int, error) {
	ctx := context.Background()
	
	// Get a Tx for making transaction requests.
    tx, err := model.DB.BeginTx(ctx, nil)
    if err != nil {
        return 0, err
    }
    // Defer a rollback in case anything fails.
    defer tx.Rollback()

	insertStmt := `INSERT INTO pool_option (pool_id, name, description) VALUES (?, ?, ?)`
	updateStmt := `UPDATE pool_option 
				   SET name = ?, description = ? 
				   WHERE option_id = ?`

	for _, option := range pool.PoolOptions {
		if option.ID == 0 {
			_, err = tx.ExecContext(ctx, insertStmt,
			pool.ID, option.Option, option.Description)
			if err != nil {
				return 0, err
			}
		} else {
			_, err = tx.ExecContext(ctx, updateStmt,
			option.Option, option.Description, option.ID)
			if err != nil {
				return 0, err
			}
		}
	}

    // Commit the transaction.
    if err = tx.Commit(); err != nil {
        return 0, err
    }

    return 0, nil
}

func (model *PoolModel) InsertConfig(pool models.Pool) (int, error) {
	stmt := `INSERT INTO pool_config (pool_id, single_vote, start_date, end_date) VALUES (?, ?, ?, ?)`
	result, err := model.DB.Exec(stmt, pool.ID, pool.PoolConfig.SingleVote, pool.PoolConfig.StartDate.UTC().Format("2006-01-02T15:04"), pool.PoolConfig.EndDate.UTC().Format("2006-01-02T15:04"))
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	} 

	return int(id), nil
}

func (model *PoolModel) GetConfig(pool models.Pool) (models.PoolConfig, error) {
	stmt := `SELECT pool_id, single_vote, start_date, end_date FROM pool_config where pool_id = ?`
	var poolConfig models.PoolConfig

	rows, err := model.DB.Query(stmt, pool.ID)
	if err != nil {
		return poolConfig, err
	}

	defer rows.Close()
	
	err = model.DB.QueryRow(stmt, pool.ID).Scan(&poolConfig.PoolID, &poolConfig.SingleVote, &poolConfig.StartDate, &poolConfig.EndDate)
	if err == sql.ErrNoRows {
		return poolConfig, models.ErrNoRecord
	} else if err != nil {
		return poolConfig, err
	}

	if err = rows.Err(); err != nil {
		return poolConfig, err
	}

	return poolConfig, nil
}

func (model *PoolModel) UpdateConfig(pool models.Pool) (int, error) {
	updateStmt := `UPDATE pool_config 
				   SET single_vote = ?, start_date = ?, end_date = ?
				   WHERE pool_id = ?`
	

	_, err := model.DB.Exec(updateStmt, pool.PoolConfig.SingleVote, pool.PoolConfig.StartDate.UTC().Format("2006-01-02T15:04"), pool.PoolConfig.EndDate.UTC().Format("2006-01-02T15:04"), pool.ID)
	if err != nil {
		return 0, err
	}

    return 0, nil
}

func (model *PoolModel) GetOptionID(id int, optionText string) (int, error) {
	stmt := `SELECT option_id FROM pool_option WHERE pool_id = ? AND name = ?`
	var optionID int

	err := model.DB.QueryRow(stmt, id, optionText).Scan(&optionID)
	if err == sql.ErrNoRows {
		return 0, models.ErrNoRecord
	} else if err != nil {
		return 0, err
	}

	return optionID, nil
}
