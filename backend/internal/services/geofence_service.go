package services

import (
	"fmt"
	"geofence/internal/domain"
	"github.com/google/uuid"
)

// GeofenceService implements domain.GeofenceService
type GeofenceService struct {
	repo domain.GeofenceRepository
}

// NewGeofenceService creates a new geofence service
func NewGeofenceService(repo domain.GeofenceRepository) domain.GeofenceService {
	return &GeofenceService{repo: repo}
}

// CreateGeofence creates a new geofence
func (s *GeofenceService) CreateGeofence(geofence *domain.Geofence) (*domain.Geofence, error) {
	// Validate coordinates
	if err := s.ValidateGeofenceCoordinates(geofence.Coordinates); err != nil {
		return nil, err
	}

	// Generate ID if not provided
	if geofence.ID == "" {
		geofence.ID = "geo_" + uuid.New().String()[:12]
	}

	// Set default status
	if geofence.Status == "" {
		geofence.Status = "active"
	}

	// Create in repository
	if err := s.repo.Create(geofence); err != nil {
		return nil, fmt.Errorf("failed to create geofence: %w", err)
	}

	return geofence, nil
}

// GetGeofences retrieves all geofences
func (s *GeofenceService) GetGeofences(category *domain.GeofenceCategory) ([]*domain.Geofence, error) {
	return s.repo.GetAll(category)
}

// GetGeofenceByID retrieves a geofence by ID
func (s *GeofenceService) GetGeofenceByID(id string) (*domain.Geofence, error) {
	return s.repo.GetByID(id)
}

// ValidateGeofenceCoordinates validates geofence coordinates
func (s *GeofenceService) ValidateGeofenceCoordinates(coordinates [][]float64) error {
	// Check minimum points (3 unique + 1 closing = 4)
	if len(coordinates) < 4 {
		return fmt.Errorf("minimum 4 coordinate points required (3 unique + 1 closing): %w", domain.ErrInvalidCoords)
	}

	// Check if first and last are identical (closed polygon)
	if len(coordinates[0]) != 2 || len(coordinates[len(coordinates)-1]) != 2 {
		return fmt.Errorf("coordinates must be [lat, lon] pairs: %w", domain.ErrInvalidCoords)
	}

	if coordinates[0][0] != coordinates[len(coordinates)-1][0] ||
		coordinates[0][1] != coordinates[len(coordinates)-1][1] {
		return fmt.Errorf("polygon must be closed (first and last coordinates must match): %w", domain.ErrInvalidCoords)
	}

	// Validate each coordinate
	for i, coord := range coordinates {
		if len(coord) != 2 {
			return fmt.Errorf("coordinate %d must have exactly 2 values [lat, lon]: %w", i, domain.ErrInvalidCoords)
		}

		lat, lon := coord[0], coord[1]

		// Validate latitude range
		if lat < -90 || lat > 90 {
			return fmt.Errorf("coordinate %d: latitude %f out of range [-90, 90]: %w", i, lat, domain.ErrInvalidCoords)
		}

		// Validate longitude range
		if lon < -180 || lon > 180 {
			return fmt.Errorf("coordinate %d: longitude %f out of range [-180, 180]: %w", i, lon, domain.ErrInvalidCoords)
		}
	}

	return nil
}

// IsPointInPolygon checks if a point is inside a polygon using ray casting algorithm
func (s *GeofenceService) IsPointInPolygon(latitude, longitude float64, geofence *domain.Geofence) bool {
	if geofence == nil || len(geofence.Coordinates) < 3 {
		return false
	}

	return s.rayIntersectsSegment(latitude, longitude, geofence.Coordinates)
}

// rayIntersectsSegment implements ray casting algorithm for point-in-polygon detection
func (s *GeofenceService) rayIntersectsSegment(lat, lon float64, polygon [][]float64) bool {
	inside := false

	for i, j := 0, len(polygon)-1; i < len(polygon); j, i = i, i+1 {
		xi, yi := polygon[i][0], polygon[i][1]
		xj, yj := polygon[j][0], polygon[j][1]

		// Check if point's latitude is between the y-coordinates of the segment
		if (yi > lat) != (yj > lat) &&
			lon < (xj-xi)*(lat-yi)/(yj-yi)+xi {
			inside = !inside
		}
	}

	return inside
}
