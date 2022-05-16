package model

type EndpointStatus int

const (
	DOWN EndpointStatus = iota
	TRANSITIONALLY_DOWN
	TRANSITIONALLY_UP
	UP
)

type EndpointType string

const (
	HTTP_ENDPOINT EndpointType = "http"
	GRPC                       = "grpc"
)

type Endpoint struct {
	Name         string         `json:"name,omitempty"`
	Target       string         `json:"target,omitempty"`
	Weight       int            `json:"weight,omitempty"`
	Backup       bool           `json:"backup,omitempty"`
	Status       EndpointStatus `json:"-,omitempty"`
	Tenants      []string       `json:"tenants,omitempty"`
	EndpointType EndpointType   `json:"type,omitempty"`
	Inherit      bool           `json:"inherit,omitempty"`
}

type EndpointHealthCheckService struct {
	Inherit bool `json:"inherit,omitempty"`

	// HealthCheckService
	Steps    []Step `json:"steps,omitempty"`
	Schedule string `json:"schedule,omitempty"`
}

func NewEndpointHealthCheckService() *EndpointHealthCheckService {
	return &EndpointHealthCheckService{Schedule: "health-check"}
}

type HttpEndpoint struct {
	// From Endpoint
	Name         string         `json:"name,omitempty"`
	Target       string         `json:"target,omitempty"`
	Weight       int            `json:"weight,omitempty"`
	Backup       bool           `json:"backup,omitempty"`
	Status       EndpointStatus `json:"-,omitempty"`
	Tenants      []string       `json:"tenants,omitempty"`
	EndpointType EndpointType   `json:"type,omitempty"`
	Inherit      bool           `json:"inherit,omitempty"`

	HttpProxy            *HttpProxy                  `json:"httpProxy,omitempty"`
	HttpClientOptions    *HttpClientOptions          `json:"httpClientOptions,omitempty"`
	HttpClientSslOptions *HttpClientSslOptions       `json:"httpClientSslOptions,omitempty"`
	Headers              map[string]string           `json:"headers,omitempty"`
	HealthCheck          *EndpointHealthCheckService `json:"healthCheck,omitempty"`
}

type EndpointDiscoveryService struct {
	Name          string            `json:"name,omitempty"`
	Enabled       bool              `json:"enabled,omitempty"`
	Service       *Service          `json:"-,omitempty"`
	Provider      string            `json:"provider,omitempty"`
	Configuration *GenericStringMap `json:"configuration,omitempty"`
}

type DynamicPropertyProvider int

const (
	HTTP_PROPERTY_PROVIDER DynamicPropertyProvider = iota
)

type DynamicPropertyService struct {
	Schedule string                  `json:"schedule,omitempty"`
	Provider DynamicPropertyProvider `json:"provider,omitempty"`
	//Configuration DynamicPropertyProviderConfiguration `json:"configuration,omitempty"`  // needs to be fixed later
}

func NewDynamicPropertyService() *DynamicPropertyService {
	return &DynamicPropertyService{
		Schedule: "dynamic-property",
	}
}

type Service struct {
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

func NewService() *Service {
	return &Service{
		Enabled: true,
	}
}

type Services struct {
	Services                 map[Service]*Service      `json:"-"`
	EndpointDiscoveryService *EndpointDiscoveryService `json:"discovery,omitempty"`
	HealthCheckService       *HealthCheckService       `json:"health-check,omitempty"`
	DynamicPropertyService   *DynamicPropertyService   `json:"dynamic-property,omitempty"`
}

type LoadBalancerType string

const (
	ROUND_ROBIN          LoadBalancerType = "ROUND_ROBIN"
	RANDOM                                = "RANDOM"
	WEIGHTED_ROUND_ROBIN                  = "WEIGHTED_ROUND_ROBIN"
	WEIGHTED_RANDOM                       = "WEIGHTED_RANDOM"
)

type LoadBalancer struct {
	LoadBalancerType LoadBalancerType `json:"type,omitempty"`
}

type EndpointGroup struct {
	Name                 string                `json:"name,omitempty"`
	Endpoints            []*HttpEndpoint       `json:"endpoints,omitempty"`
	LoadBalancer         LoadBalancer          `json:"load_balancing,omitempty"`
	Services             *Services             `json:"services,omitempty"`
	HttpProxy            *HttpProxy            `json:"proxy,omitempty"`
	HttpClientOptions    *HttpClientOptions    `json:"http,omitempty"`
	HttpClientSslOptions *HttpClientSslOptions `json:"ssl,omitempty"`
	Headers              map[string]string     `json:"headers,omitempty"`
}
