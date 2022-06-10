package mysql

import (
	"database/sql"

	"github.com/mihai22125/goPool/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// Define a PoolModel type which wraps a sql.DB connection pool
type ApiKeyModel struct {
	DB *sql.DB
}

// Insert will insert a new pool int the database
func (model *ApiKeyModel) Insert(machineID int, hashedKey []byte) (int, error) {
	stmt := `INSERT INTO api_key (machine_id, hashed_key) VALUES (?, ?)`
	result, err := model.DB.Exec(stmt, machineID, hashedKey)
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

func (m *ApiKeyModel) ValidateKey(machineID int, key string) (int, error) {
	var id int
	var hashedKey []byte
	row := m.DB.QueryRow("SELECT machine_id, hashed_key FROM api_key WHERE machine_id = ?", machineID)
	err := row.Scan(&id, &hashedKey)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedKey, []byte(key))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	return id, nil
}
