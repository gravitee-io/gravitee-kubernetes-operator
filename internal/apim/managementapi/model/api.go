package model

type Api struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	State             string `json:"state"`
	Visibility        string `json:"visibility"`
	ApiLifecycleState string `json:"lifecycle_state"`
}

type Action string

const (
	ActionStart Action = "START"
	ActionStop  Action = "STOP"
)
