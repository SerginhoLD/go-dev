package entity

import "time"

type Object struct {
	Id        uint64
	Title     string
	Metro     string
	Price     uint64
	Size      float64
	Rooms     uint8
	Checked   bool
	UpdatedAt time.Time
}
