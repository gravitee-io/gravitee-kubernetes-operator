package model

type Api struct {
	Id                string  `json:"id"`
	Name              string  `json:"name"`
	State             string  `json:"state"`
	Visibility        string  `json:"visibility"`
	ApiLifecycleState string  `json:"lifecycle_state"`
	Plans             []*Plan `json:"plans"`
}

type Action string

const (
	ActionStart Action = "START"
	ActionStop  Action = "STOP"
)

type Plan struct {
	Id       string           `json:"id"`
	CrossId  string           `json:"crossId"`
	Name     string           `json:"name"`
	Security PlanSecurityType `json:"security"`
	Status   PlanStatus       `json:"status"`
}

type PlanSecurityType string

const (
	PlanSecurityTypeKeyLess PlanSecurityType = "KEY_LESS"
	PlanSecurityTypeApiKey  PlanSecurityType = "API_KEY"
	PlanSecurityTypeOauth2  PlanSecurityType = "OAUTH2"
	PlanSecurityTypeJwt     PlanSecurityType = "JWT"
)

type PlanStatus string

const (
	PlanStatusStaging    PlanStatus = "STAGING"
	PlanStatusPublished  PlanStatus = "PUBLISHED"
	PlanStatusDeprecated PlanStatus = "DEPRECATED"
	PlanStatusClosed     PlanStatus = "CLOSED"
)
