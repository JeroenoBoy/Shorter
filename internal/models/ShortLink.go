package models

import (
	"strconv"
	"time"
)

type LinkId int

type ShortLink struct {
	Id        LinkId     `json:"id"`
	Owner     UserId     `json:"owner"`
	Link     string     `json:"link"`
	Target    string     `json:"target"`
	Redirects int        `json:"redirects"`
	CreatedAt time.Time  `json:"createdAt"`
	LastUsed  *time.Time `json:"lastUsed"`
}

func (l LinkId) ToString() string {
	return strconv.Itoa(int(l))
}
