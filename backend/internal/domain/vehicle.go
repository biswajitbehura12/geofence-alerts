package domain

import "time"

// VehicleStatus represents the status of a vehicle
type VehicleStatus string

const (
	VehicleActive   VehicleStatus = "active"
	VehicleInactive VehicleStatus = "inactive"
)

// Vehicle represents a registered vehicle in the system
type Vehicle struct {
	ID            string        `json:"id"`
	VehicleNumber string        `json:"vehicle_number"`
	DriverName    string        `json:"driver_name"`
	VehicleType   string        `json:"vehicle_type"`
	Phone         string        `json:"phone"`
	Status        VehicleStatus `json:"status"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// VehicleLocation represents the current location of a vehicle
type VehicleLocation struct {
	ID        string    `json:"id"`
	VehicleID string    `json:"vehicle_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

// VehicleRepository defines the interface for vehicle persistence
type VehicleRepository interface {
	Create(vehicle *Vehicle) error
	GetByID(id string) (*Vehicle, error)
	GetAll() ([]*Vehicle, error)
	Update(vehicle *Vehicle) error
	Delete(id string) error
	GetByVehicleNumber(vehicleNumber string) (*Vehicle, error)
}

// VehicleLocationRepository defines the interface for vehicle location persistence
type VehicleLocationRepository interface {
	SaveLocation(location *VehicleLocation) error
	GetLatestLocation(vehicleID string) (*VehicleLocation, error)
	GetLocationHistory(vehicleID string, limit int, offset int) ([]*VehicleLocation, int64, error)
}

// VehicleService defines the interface for vehicle business logic
type VehicleService interface {
	RegisterVehicle(vehicle *Vehicle) (*Vehicle, error)
	GetVehicles() ([]*Vehicle, error)
	GetVehicleByID(id string) (*Vehicle, error)
	UpdateVehicleLocation(vehicleID string, latitude, longitude float64, timestamp time.Time) (*VehicleLocation, error)
	GetVehicleLocation(vehicleID string) (*VehicleLocation, error)
	GetVehicleLocationHistory(vehicleID string, limit, offset int) ([]*VehicleLocation, int64, error)
}
