package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

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