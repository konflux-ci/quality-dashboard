package server

import "time"

// HealthCheckStatus contains specification and current status of external services
// and contains a map of currently unhealthy critical components
type HealthCheckStatus struct {
	ExternalServices            []Service           `json:"externalServices"`
	UnhealthyCriticalComponents map[string][]string `json:"unhealthyCriticalComponents"`
}

// Service contains specification of an external service specified via config
// and holds an information about current status of related components
type Service struct {
	Name               string   `json:"name"`
	CriticalComponents []string `json:"criticalComponents"`
	StatusPageURL      string   `json:"statusPageURL"`
	CurrentStatus      Summary  `json:"currentStatus"`
}

// Summary is the Statuspage API component representation
type Summary struct {
	Components []Component `json:"components"`
	Incidents  []Incident  `json:"incidents"`
	Status     Status      `json:"status"`
}

// Component is the Statuspage API component representation
type Component struct {
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Name               string
	GroupID            string
	PageID             string
	ID                 string
	Description        string
	Status             string
	AutomationEmail    string
	Position           int32
	Group              bool
	Showcase           bool
	OnlyShowIfDegraded bool
}

// Incident entity reflects one single incident
type Incident struct {
	ID                string           `json:"id,omitempty"`
	Name              string           `json:"name,omitempty"`
	Status            string           `json:"status,omitempty"`
	Message           string           `json:"message,omitempty"`
	Visible           int              `json:"visible,omitempty"`
	ComponentID       int              `json:"component_id,omitempty"`
	ComponentStatus   int              `json:"component_status,omitempty"`
	Notify            bool             `json:"notify,omitempty"`
	Stickied          bool             `json:"stickied,omitempty"`
	OccurredAt        string           `json:"occurred_at,omitempty"`
	Template          string           `json:"template,omitempty"`
	Vars              []string         `json:"vars,omitempty"`
	CreatedAt         string           `json:"created_at,omitempty"`
	UpdatedAt         string           `json:"updated_at,omitempty"`
	DeletedAt         string           `json:"deleted_at,omitempty"`
	IsResolved        bool             `json:"is_resolved,omitempty"`
	Updates           []IncidentUpdate `json:"incident_updates,omitempty"`
	HumanStatus       string           `json:"human_status,omitempty"`
	LatestUpdateID    int              `json:"latest_update_id,omitempty"`
	LatestStatus      int              `json:"latest_status,omitempty"`
	LatestHumanStatus string           `json:"latest_human_status,omitempty"`
	LatestIcon        string           `json:"latest_icon,omitempty"`
	Permalink         string           `json:"permalink,omitempty"`
	Duration          int              `json:"duration,omitempty"`
}

// IncidentUpdate entity reflects one single incident update
type IncidentUpdate struct {
	ID              string `json:"id,omitempty"`
	Body            string `json:"body,omitempty"`
	IncidentID      string `json:"incident_id,omitempty"`
	ComponentID     int    `json:"component_id,omitempty"`
	ComponentStatus int    `json:"component_status,omitempty"`
	Status          string `json:"status,omitempty"`
	Message         string `json:"message,omitempty"`
	UserID          int    `json:"user_id,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
	HumanStatus     string `json:"human_status,omitempty"`
	Permalink       string `json:"permalink,omitempty"`
}

// Status entity contains the contents of API Response of a /status call.
type Status struct {
	Indicator   string `json:"indicator,omitempty"`
	Description string `json:"description,omitempty"`
}
