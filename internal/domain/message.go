package domain

import "time"

type Message struct {
	UserID    string
	Username  string
	Content   string
	TimeStamp time.Time
}
