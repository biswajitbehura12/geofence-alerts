package models

// HTTP Request/Response models

// CreateGeofenceRequest is the request body for POST /geofences
type CreateGeofenceRequest struct {
	Name        string      `json:"name" binding:"required"`
	Description string      `json:"description"`
	Coordinates [][]float64 `json:"coordinates" binding:"required"`
	Category    string      `json:"category" binding:"required"`
}

// GeofenceResponse is the response for geofence operations
type GeofenceResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	TimeNS string `json:"time_ns"`
}

// GeofenceListResponse is the response for GET /geofences
type GeofenceListResponse struct {
	Geofences []*GeofenceDetail `json:"geofences"`
	TimeNS    string            `json:"time_ns"`
}

// GeofenceDetail contains geofence details
type GeofenceDetail struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Coordinates [][]float64 `json:"coordinates"`
	Category    string      `json:"category"`
	CreatedAt   string      `json:"created_at"`
}

// CreateVehicleRequest is the request body for POST /vehicles
type CreateVehicleRequest struct {
	VehicleNumber string `json:"vehicle_number" binding:"required"`
	DriverName    string `json:"driver_name" binding:"required"`
	VehicleType   string `json:"vehicle_type" binding:"required"`
	Phone         string `json:"phone" binding:"required"`
}

// VehicleResponse is the response for vehicle operations
type VehicleResponse struct {
	ID            string `json:"id"`
	VehicleNumber string `json:"vehicle_number"`
	Status        string `json:"status"`
	TimeNS        string `json:"time_ns"`
}

// VehicleListResponse is the response for GET /vehicles
type VehicleListResponse struct {
	Vehicles []*VehicleDetail `json:"vehicles"`
	TimeNS   string           `json:"time_ns"`
}

// VehicleDetail contains vehicle details
type VehicleDetail struct {
	ID            string `json:"id"`
	VehicleNumber string `json:"vehicle_number"`
	DriverName    string `json:"driver_name"`
	VehicleType   string `json:"vehicle_type"`
	Phone         string `json:"phone"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

// UpdateLocationRequest is the request body for POST /vehicles/location
type UpdateLocationRequest struct {
	VehicleID string `json:"vehicle_id" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Timestamp string  `json:"timestamp" binding:"required"`
}

// LocationUpdateResponse is the response for POST /vehicles/location
type LocationUpdateResponse struct {
	VehicleID        string                    `json:"vehicle_id"`
	LocationUpdated  bool                      `json:"location_updated"`
	CurrentGeofences []CurrentGeofenceStatus   `json:"current_geofences"`
	TimeNS           string                    `json:"time_ns"`
}

// CurrentGeofenceStatus contains current geofence status
type CurrentGeofenceStatus struct {
	GeofenceID   string `json:"geofence_id"`
	GeofenceName string `json:"geofence_name"`
	Status       string `json:"status"`
}

// VehicleLocationResponse is the response for GET /vehicles/location/{vehicle_id}
type VehicleLocationResponse struct {
	VehicleID        string                    `json:"vehicle_id"`
	VehicleNumber    string                    `json:"vehicle_number"`
	CurrentLocation  Location                  `json:"current_location"`
	CurrentGeofences []CurrentGeofenceStatus   `json:"current_geofences"`
	TimeNS           string                    `json:"time_ns"`
}

// Location contains location details
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}

// ConfigureAlertRequest is the request body for POST /alerts/configure
type ConfigureAlertRequest struct {
	GeofenceID string `json:"geofence_id" binding:"required"`
	VehicleID  string `json:"vehicle_id"`
	EventType  string `json:"event_type" binding:"required"`
}

// AlertResponse is the response for alert operations
type AlertResponse struct {
	AlertID    string `json:"alert_id"`
	GeofenceID string `json:"geofence_id"`
	VehicleID  string `json:"vehicle_id"`
	EventType  string `json:"event_type"`
	Status     string `json:"status"`
	TimeNS     string `json:"time_ns"`
}

// AlertListResponse is the response for GET /alerts
type AlertListResponse struct {
	Alerts []*AlertDetail `json:"alerts"`
	TimeNS string         `json:"time_ns"`
}

// AlertDetail contains alert details
type AlertDetail struct {
	AlertID       string `json:"alert_id"`
	GeofenceID    string `json:"geofence_id"`
	GeofenceName  string `json:"geofence_name"`
	VehicleID     string `json:"vehicle_id"`
	VehicleNumber string `json:"vehicle_number"`
	EventType     string `json:"event_type"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

// ViolationHistoryResponse is the response for GET /violations/history
type ViolationHistoryResponse struct {
	Violations []*ViolationDetail `json:"violations"`
	TotalCount int64              `json:"total_count"`
	TimeNS     string             `json:"time_ns"`
}

// ViolationDetail contains violation details
type ViolationDetail struct {
	ID            string  `json:"id"`
	VehicleID     string  `json:"vehicle_id"`
	VehicleNumber string  `json:"vehicle_number"`
	GeofenceID    string  `json:"geofence_id"`
	GeofenceName  string  `json:"geofence_name"`
	EventType     string  `json:"event_type"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Timestamp     string  `json:"timestamp"`
}

// ErrorResponse is the standard error response
type ErrorResponse struct {
	Error      string `json:"error"`
	Code       string `json:"code,omitempty"`
	TimeNS     string `json:"time_ns"`
}
