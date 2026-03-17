package domain

import "time"

// Coordinate represents a latitude/longitude pair
type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// GeofenceCategory defines the type of geofence
type GeofenceCategory string

const (
	DeliveryZone   GeofenceCategory = "delivery_zone"
	RestrictedZone GeofenceCategory = "restricted_zone"
	TollZone       GeofenceCategory = "toll_zone"
	CustomerArea   GeofenceCategory = "customer_area"
)

// Geofence represents a virtual boundary
type Geofence struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Coordinates [][]float64         `json:"coordinates"`
	Category    GeofenceCategory    `json:"category"`
	Status      string              `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// GeofenceRepository defines the interface for geofence persistence
type GeofenceRepository interface {
	Create(geofence *Geofence) error
	GetByID(id string) (*Geofence, error)
	GetAll(category *GeofenceCategory) ([]*Geofence, error)
	Update(geofence *Geofence) error
	Delete(id string) error
}

// GeofenceService defines the interface for geofence business logic
type GeofenceService interface {
	CreateGeofence(geofence *Geofence) (*Geofence, error)
	GetGeofences(category *GeofenceCategory) ([]*Geofence, error)
	GetGeofenceByID(id string) (*Geofence, error)
	ValidateGeofenceCoordinates(coordinates [][]float64) error
	IsPointInPolygon(latitude, longitude float64, geofence *Geofence) bool
}
