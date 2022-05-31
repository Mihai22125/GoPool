package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail = errors.New("models: duplicate email")

	RoleAdmin = "admin"
	RoleBasic = "basic"
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

//TODO: create user role type
type User struct {
	ID             int
	Name           string
	Email          string
	Role 		   string
	HashedPassword []byte
	Created        time.Time
}

type Machine struct {
	ID           int
	PhoneNumber  string
	IPAdrres     string
}

type Session struct {
	ID        int
	MachineID int
	PoolID    int
}

type Vote struct {
	ID        int
	PoolID    int
	OptionID  int
	MachineID int
	From      string
}

type VoteRequest struct {
	MachineID int     `json:"machine_id"`
	Text 	  string  `json:"text"`
	From 	  string  `json:"from"`
}

type Result struct {
	Option     PoolOption
	Count      int
	Percentage float64
}