package models

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/lib/pq"
	"time"
)

// Geofence database model
type Geofence struct {
	ID          string           `db:"id"`
	Name        string           `db:"name"`
	Description string           `db:"description"`
	Coordinates pq.Float64Array  `db:"coordinates"` // Flattened array: [lat1, lon1, lat2, lon2, ...]
	Category    string           `db:"category"`
	Status      string           `db:"status"`
	CreatedAt   time.Time        `db:"created_at"`
	UpdatedAt   time.Time        `db:"updated_at"`
}

// Vehicle database model
type Vehicle struct {
	ID            string    `db:"id"`
	VehicleNumber string    `db:"vehicle_number"`
	DriverName    string    `db:"driver_name"`
	VehicleType   string    `db:"vehicle_type"`
	Phone         string    `db:"phone"`
	Status        string    `db:"status"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// VehicleLocation database model
type VehicleLocation struct {
	ID        string    `db:"id"`
	VehicleID string    `db:"vehicle_id"`
	Latitude  float64   `db:"latitude"`
	Longitude float64   `db:"longitude"`
	Timestamp time.Time `db:"timestamp"`
	CreatedAt time.Time `db:"created_at"`
}

// AlertRule database model
type AlertRule struct {
	ID        string    `db:"id"`
	GeofenceID string    `db:"geofence_id"`
	VehicleID string    `db:"vehicle_id"`
	EventType string    `db:"event_type"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// GeofenceEvent database model
type GeofenceEvent struct {
	ID         string    `db:"id"`
	VehicleID  string    `db:"vehicle_id"`
	GeofenceID string    `db:"geofence_id"`
	EventType  string    `db:"event_type"`
	Latitude   float64   `db:"latitude"`
	Longitude  float64   `db:"longitude"`
	Timestamp  time.Time `db:"timestamp"`
	CreatedAt  time.Time `db:"created_at"`
}

// GeofenceWithNames is used for API responses
type GeofenceWithNames struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// VehicleWithStatus is used for API responses
type VehicleWithStatus struct {
	ID            string    `json:"id"`
	VehicleNumber string    `json:"vehicle_number"`
	DriverName    string    `json:"driver_name"`
	VehicleType   string    `json:"vehicle_type"`
	Phone         string    `json:"phone"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// VehicleLocationInfo contains location and geofence info
type VehicleLocationInfo struct {
	CurrentLocation CurrentLocation `json:"current_location"`
	CurrentGeofences []GeofenceInfo  `json:"current_geofences"`
}

// CurrentLocation contains location details
type CurrentLocation struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
}

// GeofenceInfo contains geofence details for location response
type GeofenceInfo struct {
	GeofenceID   string `json:"geofence_id"`
	GeofenceName string `json:"geofence_name"`
	Category     string `json:"category"`
}

// AlertRuleWithNames is used for API responses
type AlertRuleWithNames struct {
	AlertID       string    `json:"alert_id"`
	GeofenceID    string    `json:"geofence_id"`
	GeofenceName  string    `json:"geofence_name"`
	VehicleID     string    `json:"vehicle_id"`
	VehicleNumber string    `json:"vehicle_number"`
	EventType     string    `json:"event_type"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// GeofenceEventWithNames is used for API responses
type GeofenceEventWithNames struct {
	ID               string    `json:"id"`
	VehicleID        string    `json:"vehicle_id"`
	VehicleNumber    string    `json:"vehicle_number"`
	GeofenceID       string    `json:"geofence_id"`
	GeofenceName     string    `json:"geofence_name"`
	EventType        string    `json:"event_type"`
	Latitude         float64   `json:"latitude"`
	Longitude        float64   `json:"longitude"`
	Timestamp        time.Time `json:"timestamp"`
}

// LocationCoordinates for JSON marshaling
type LocationCoordinates [][]float64

// Scan implements sql.Scanner
func (lc *LocationCoordinates) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, _ := value.([]byte)
	return json.Unmarshal(b, &lc)
}

// Value implements driver.Valuer
func (lc LocationCoordinates) Value() (driver.Value, error) {
	return json.Marshal(lc)
}
