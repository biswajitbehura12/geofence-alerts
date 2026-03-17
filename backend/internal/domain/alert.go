package domain

import "time"

// EventType represents the type of geofence event
type EventType string

const (
	EventEntry EventType = "entry"
	EventExit  EventType = "exit"
	EventBoth  EventType = "both"
)

// AlertRule represents a rule for triggering alerts
type AlertRule struct {
	ID          string    `json:"alert_id"`
	GeofenceID  string    `json:"geofence_id"`
	VehicleID   string    `json:"vehicle_id"` // Empty string means all vehicles
	EventType   EventType `json:"event_type"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GeofenceEvent represents a violation event
type GeofenceEvent struct {
	ID         string    `json:"id"`
	VehicleID  string    `json:"vehicle_id"`
	GeofenceID string    `json:"geofence_id"`
	EventType  EventType `json:"event_type"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Timestamp  time.Time `json:"timestamp"`
	CreatedAt  time.Time `json:"created_at"`
}

// RealTimeAlert represents an alert sent to clients
type RealTimeAlert struct {
	EventID   string                 `json:"event_id"`
	EventType EventType              `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	Vehicle   VehicleAlertInfo       `json:"vehicle"`
	Geofence  GeofenceAlertInfo      `json:"geofence"`
	Location  LocationAlertInfo      `json:"location"`
}

// VehicleAlertInfo contains vehicle info for alerts
type VehicleAlertInfo struct {
	VehicleID    string `json:"vehicle_id"`
	VehicleNumber string `json:"vehicle_number"`
	DriverName   string `json:"driver_name"`
}

// GeofenceAlertInfo contains geofence info for alerts
type GeofenceAlertInfo struct {
	GeofenceID   string           `json:"geofence_id"`
	GeofenceName string           `json:"geofence_name"`
	Category    GeofenceCategory `json:"category"`
}

// LocationAlertInfo contains location info for alerts
type LocationAlertInfo struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// AlertRepository defines the interface for alert persistence
type AlertRepository interface {
	CreateRule(rule *AlertRule) error
	GetRuleByID(id string) (*AlertRule, error)
	GetAllRules(geofenceID, vehicleID *string) ([]*AlertRule, error)
	DeleteRule(id string) error
	SaveEvent(event *GeofenceEvent) error
	GetEventHistory(vehicleID, geofenceID *string, startDate, endDate *time.Time, limit, offset int) ([]*GeofenceEvent, int64, error)
}

// AlertService defines the interface for alert business logic
type AlertService interface {
	ConfigureAlert(geofenceID, vehicleID string, eventType EventType) (*AlertRule, error)
	GetAlerts(geofenceID, vehicleID *string) ([]*AlertRule, error)
	GetViolationHistory(vehicleID, geofenceID *string, startDate, endDate *time.Time, limit, offset int) ([]*GeofenceEvent, int64, error)
	CheckAndTriggerAlerts(vehicle *Vehicle, geofences []*Geofence, eventType EventType) ([]*RealTimeAlert, error)
}

// AlertPublisher defines the interface for publishing alerts
type AlertPublisher interface {
	PublishAlert(alert *RealTimeAlert) error
	Subscribe(ch chan *RealTimeAlert) string
	Unsubscribe(subscriptionID string)
}
