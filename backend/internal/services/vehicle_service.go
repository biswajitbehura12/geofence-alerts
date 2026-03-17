package services

import (
	"fmt"
	"geofence/internal/domain"
	"github.com/google/uuid"
	"time"
)

// VehicleService implements domain.VehicleService
type VehicleService struct {
	vehicleRepo          domain.VehicleRepository
	locationRepo         domain.VehicleLocationRepository
	geofenceRepo         domain.GeofenceRepository
	geofenceService      domain.GeofenceService
}

// NewVehicleService creates a new vehicle service
func NewVehicleService(
	vehicleRepo domain.VehicleRepository,
	locationRepo domain.VehicleLocationRepository,
	geofenceRepo domain.GeofenceRepository,
	geofenceService domain.GeofenceService,
) domain.VehicleService {
	return &VehicleService{
		vehicleRepo:     vehicleRepo,
		locationRepo:    locationRepo,
		geofenceRepo:    geofenceRepo,
		geofenceService: geofenceService,
	}
}

// RegisterVehicle registers a new vehicle
func (s *VehicleService) RegisterVehicle(vehicle *domain.Vehicle) (*domain.Vehicle, error) {
	// Check if vehicle already exists
	existing, _ := s.vehicleRepo.GetByVehicleNumber(vehicle.VehicleNumber)
	if existing != nil {
		return nil, domain.ErrVehicleExists
	}

	// Generate ID if not provided
	if vehicle.ID == "" {
		vehicle.ID = "veh_" + uuid.New().String()[:12]
	}

	// Set default status
	if vehicle.Status == "" {
		vehicle.Status = domain.VehicleActive
	}

	// Create in repository
	if err := s.vehicleRepo.Create(vehicle); err != nil {
		return nil, fmt.Errorf("failed to register vehicle: %w", err)
	}

	return vehicle, nil
}

// GetVehicles retrieves all vehicles
func (s *VehicleService) GetVehicles() ([]*domain.Vehicle, error) {
	return s.vehicleRepo.GetAll()
}

// GetVehicleByID retrieves a vehicle by ID
func (s *VehicleService) GetVehicleByID(id string) (*domain.Vehicle, error) {
	return s.vehicleRepo.GetByID(id)
}

// UpdateVehicleLocation updates a vehicle's location and detects geofence entry/exit
func (s *VehicleService) UpdateVehicleLocation(vehicleID string, latitude, longitude float64, timestamp time.Time) (*domain.VehicleLocation, error) {
	// Validate vehicle exists
	_, err := s.vehicleRepo.GetByID(vehicleID)
	if err != nil {
		return nil, err
	}

	// Create location record
	location := &domain.VehicleLocation{
		ID:        "loc_" + uuid.New().String()[:12],
		VehicleID: vehicleID,
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: timestamp,
	}

	// Save location
	if err := s.locationRepo.SaveLocation(location); err != nil {
		return nil, fmt.Errorf("failed to save location: %w", err)
	}

	return location, nil
}

// GetVehicleLocation gets the current location of a vehicle
func (s *VehicleService) GetVehicleLocation(vehicleID string) (*domain.VehicleLocation, error) {
	// Validate vehicle exists
	_ , err := s.vehicleRepo.GetByID(vehicleID)
	if err != nil {
		return nil, err
	}

	location, err := s.locationRepo.GetLatestLocation(vehicleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}

	return location, nil
}

// GetVehicleLocationHistory gets location history for a vehicle
func (s *VehicleService) GetVehicleLocationHistory(vehicleID string, limit, offset int) ([]*domain.VehicleLocation, int64, error) {
	// Validate vehicle exists
	_ , err := s.vehicleRepo.GetByID(vehicleID)
	if err != nil {
		return nil, 0, err
	}

	// Default limit
	if limit == 0 {
		limit = 50
	}

	// Max limit
	if limit > 500 {
		limit = 500
	}

	return s.locationRepo.GetLocationHistory(vehicleID, limit, offset)
}

// GetCurrentGeofences returns all geofences containing the vehicle
func (s *VehicleService) GetCurrentGeofences(vehicleID string) ([]*domain.Geofence, error) {
	// Get current location
	location, err := s.GetVehicleLocation(vehicleID)
	if err != nil || location == nil {
		return nil, nil
	}

	// Get all geofences
	geofences, err := s.geofenceRepo.GetAll(nil)
	if err != nil {
		return nil, err
	}

	// Filter geofences containing the vehicle
	var current []*domain.Geofence
	for _, geofence := range geofences {
		if s.geofenceService.IsPointInPolygon(location.Latitude, location.Longitude, geofence) {
			current = append(current, geofence)
		}
	}

	return current, nil
}
