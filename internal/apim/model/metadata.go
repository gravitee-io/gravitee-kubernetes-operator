package model

type Metadata interface {
	GetName() string
	GetValue() string
	GetDefaultValue() string
}

type BaseMetadata struct {
	Name         string  `json:"name"`
	Value        *string `json:"value,omitempty"`
	DefaultValue *string `json:"defaultValue,omitempty"`
}

func (m BaseMetadata) GetName() string {
	return m.Name
}

func (m BaseMetadata) GetValue() string {
	return *m.Value
}

func (m BaseMetadata) GetDefaultValue() string {
	return *m.DefaultValue
}
