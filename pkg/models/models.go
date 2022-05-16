package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

//TODO: make ID fields int64

type Pool struct {
	ID              int
	UserID	        int
	Name            string
	NumberOfOptions int
	PoolConfig      PoolConfig
	PoolOptions     []PoolOption
}

type PoolConfig struct {
	PoolID     int
	SingleVote bool
	StartDate  time.Time
	EndDate    time.Time
}

type PoolOption struct {
	ID          int
	PoolID      int
	Option      string
	Description string
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}