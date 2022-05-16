package model

type Policy struct {
	Name          string            `json:"name,omitempty"`
	Configuration *GenericStringMap `json:"configuration,omitempty"`
}

type Rule struct {
	Methods     []HttpMethod `json:"methods,omitempty"`
	Policy      *Policy      `json:"policy,omitempty"`
	Description string       `json:"description,omitempty"`
	Enabled     bool         `json:"enabled,omitempty"`
}

func (r *Rule) ToMap() map[string]interface{} {
	m := make(map[string]interface{}, 1)
	m["enabled"] = r.Enabled
	m["description"] = r.Description
	m["methods"] = r.Methods

	if r.Policy != nil {
		m[r.Policy.Name] = r.Policy.Configuration
	}

	return m
}

func NewRule() *Rule {
	return &Rule{
		Methods: []HttpMethod{"CONNECT",
			"DELETE",
			"GET",
			"HEAD",
			"OPTIONS",
			"PATCH",
			"POST",
			"PUT",
			"TRACE",
			"OTHER"},
		Enabled: true,
	}
}

// +kubebuilder:validation:Enum=STARTS_WITH;EQUALS;
type Operator string

type PathOperator struct {
	Path string `json:"path,omitempty"`
	// +kubebuilder:default:=STARTS_WITH
	Operator Operator `json:"operator,omitempty"`
}

func NewPathOperator() *PathOperator {
	return &PathOperator{
		Operator: "STARTS_WITH",
	}
}

type FlowStep struct {
	Name          string            `json:"name,omitempty"`
	Policy        string            `json:"policy,omitempty"`
	Description   string            `json:"description,omitempty"`
	Configuration *GenericStringMap `json:"configuration,omitempty"`
	Enabled       bool              `json:"enabled,omitempty"`
	Condition     string            `json:"condition,omitempty"`
}

func NewFlowStep() *FlowStep {
	return &FlowStep{
		Enabled: true,
	}
}

type ConsumerType int

const (
	TAG ConsumerType = iota
)

type Consumer struct {
	ConsumerType ConsumerType `json:"consumerType,omitempty"`
	ConsumerId   string       `json:"consumerId,omitempty"`
}

type Flow struct {
	Name         string        `json:"name,omitempty"`
	PathOperator *PathOperator `json:"path-operator,omitempty"`
	Pre          []FlowStep    `json:"pre,omitempty"`
	Post         []FlowStep    `json:"post,omitempty"`
	Enabled      bool          `json:"enabled,omitempty"`
	Methods      []HttpMethod  `json:"methods,omitempty"`
	Condition    string        `json:"condition,omitempty"`
	Consumers    []Consumer    `json:"consumers,omitempty"`
}

func NewFlow() *Flow {
	return &Flow{
		PathOperator: NewPathOperator(),
		Enabled:      true,
	}
}

type Property struct {
	Key       string `json:"key,omitempty"`
	Value     string `json:"value,omitempty"`
	Encrypted bool   `json:"encrypted,omitempty"`
}

type ResponseTemplate struct {
	StatusCode int               `json:"status,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

// +kubebuilder:validation:Enum=AUTO;MANUAL;
type PlanValidation string

// +kubebuilder:validation:Enum=API;CATALOG;
type PlanType string

// +kubebuilder:validation:Enum=STAGING;PUBLISHED;CLOSED;DEPRECATED;
type PlanStatus string

type Plan struct {
	Id                 string            `json:"id,omitempty"`
	CrossId            string            `json:"crossId,omitempty"`
	Name               string            `json:"name"`
	Description        string            `json:"description"`
	Security           string            `json:"security"`
	SecurityDefinition string            `json:"securityDefinition,omitempty"`
	Paths              map[string][]Rule `json:"paths,omitempty"`
	Api                string            `json:"api,omitempty"`
	SelectionRule      string            `json:"selectionRule,omitempty"`
	Flows              []Flow            `json:"flows,omitempty"`
	Tags               []string          `json:"tags,omitempty"`
	// +kubebuilder:default:=PUBLISHED
	Status          PlanStatus `json:"status,omitempty"`
	Characteristics []string   `json:"characteristics,omitempty"`
	// +kubebuilder:default:=AUTO
	Validation      PlanValidation `json:"validation,omitempty"`
	CommentRequired bool           `json:"comment_required,omitempty"`
	Order           int            `json:"order,omitempty"`
	// +kubebuilder:default:=API
	Type PlanType `json:"type,omitempty"`
}

type Path struct {
	Path  string  `json:"path,omitempty"`
	Rules []*Rule `json:"rules,omitempty"`
}
