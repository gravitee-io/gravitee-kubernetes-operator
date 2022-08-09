package model

type Api struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

type Action string

const (
	ActionStart Action = "START"
	ActionStop  Action = "STOP"
)
