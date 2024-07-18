package models

type ShortId string

type ShortLink struct {
	Id     ShortId `json:"id"`
	Owner  UserId  `json:"owner"`
	Target string  `json:"target"`
}
