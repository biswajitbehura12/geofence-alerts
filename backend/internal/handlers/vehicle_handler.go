package handlers

import (
	"geofence/internal/domain"
	"geofence/internal/models"
	"geofence/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// VehicleHandler handles vehicle-related HTTP requests
type VehicleHandler struct {
	vehicleService domain.VehicleService
	alertService   domain.AlertService
	geofenceService domain.GeofenceService
}

// NewVehicleHandler creates a new vehicle handler
func NewVehicleHandler(
	vehicleService domain.VehicleService,
	alertService domain.AlertService,
	geofenceService domain.GeofenceService,
) *VehicleHandler {
	return &VehicleHandler{
		vehicleService: vehicleService,
		alertService:   alertService,
		geofenceService: geofenceService,
	}
}

// RegisterVehicle handles POST /vehicles
func (h *VehicleHandler) RegisterVehicle(c *gin.Context) {
	timer := services.NewTimeHelper()

	var req models.CreateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:  "Invalid request body",
			Code:   domain.ErrCodeInvalidInput,
			TimeNS: timer.GetElapsedNano(),
		})
		return
	}

	// Create domain vehicle
	vehicle := &domain.Vehicle{
		VehicleNumber: req.VehicleNumber,
		DriverName:    req.DriverName,
		VehicleType:   req.VehicleType,
		Phone:         req.Phone,
		Status:        domain.VehicleActive,
	}

	// Register vehicle
	result, err := h.vehicleService.RegisterVehicle(vehicle)
	if err != nil {
		code := domain.ErrCodeInternalError
		status := http.StatusInternalServerError

		if err == domain.ErrVehicleExists {
			code = domain.ErrCodeDuplicate
			status = http.StatusConflict
		}

		c.JSON(status, models.ErrorResponse{
			Error:  err.Error(),
			Code:   code,
			TimeNS: timer.GetElapsedNano(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.VehicleResponse{
		ID:            result.ID,
		VehicleNumber: result.VehicleNumber,
		Status:        string(result.Status),
		TimeNS:        timer.GetElapsedNano(),
	})
}

// GetVehicles handles GET /vehicles
func (h *VehicleHandler) GetVehicles(c *gin.Context) {
	timer := services.NewTimeHelper()

	// Get vehicles
	vehicles, err := h.vehicleService.GetVehicles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:  err.Error(),
			Code:   domain.ErrCodeInternalError,
			TimeNS: timer.GetElapsedNano(),
		})
		return
	}

	// Convert to response
	details := []*models.VehicleDetail{}
	for _, v := range vehicles {
		details = append(details, &models.VehicleDetail{
			ID:            v.ID,
			VehicleNumber: v.VehicleNumber,
			DriverName:    v.DriverName,
			VehicleType:   v.VehicleType,
			Phone:         v.Phone,
			Status:        string(v.Status),
			CreatedAt:     v.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, models.VehicleListResponse{
		Vehicles: details,
		TimeNS:   timer.GetElapsedNano(),
	})
}

// UpdateLocation handles POST /vehicles/location
func (h *VehicleHandler) UpdateLocation(c *gin.Context) {
	timer := services.NewTimeHelper()

	var req models.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:  "Invalid request body",
			Code:   domain.ErrCodeInvalidInput,
			TimeNS: timer.GetElapsedNano(),
		})
		return
	}

	// Validate coordinates
	if req.Latitude < -90 || req.Latitude > 90 || req.Longitude < -180 || req.Longitude > 180 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:  "Invalid coordinates",
			Code:   domain.ErrCodeInvalidCoords,
			TimeNS: timer.GetElapsedNano(),
		})
		return
	}

	// Parse timestamp
	timestamp, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:  "Invalid timestamp format",
			Code:   domain.ErrCodeInvalidInput,
			TimeNS: timer.GetElapsedNano(),
		})
		return
	}

	// Update location
	location, err := h.vehicleService.UpdateVehicleLocation(req.VehicleID, req.Latitude, req.Longitude, timestamp)
	if err != nil {
		code := domain.ErrCodeInternalError
		status := http.StatusInternalServerError

		if err == domain.ErrVehicleNotFound {
			code = domain.ErrCodeNotFound
			status = http.StatusNotFound
		}

		c.JSON(status, models.ErrorResponse{
			Error:  err.Error(),
			Code:   code,
			TimeNS: timer.GetElapsedNano(),
		})
		return
	}

	c.JSON(http.StatusOK, models.LocationUpdateResponse{
		VehicleID:       location.VehicleID,
		LocationUpdated: true,
		CurrentGeofences: []models.CurrentGeofenceStatus{},
		TimeNS:          timer.GetElapsedNano(),
	})
}

// GetLocation handles GET /vehicles/location/{vehicle_id}
func (h *VehicleHandler) GetLocation(c *gin.Context) {
	timer := services.NewTimeHelper()

	vehicleID := c.Param("vehicle_id")

	// Get vehicle
	vehicle, err := h.vehicleService.GetVehicleByID(vehicleID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:  "Vehicle not found",
			Code:   domain.ErrCodeNotFound,
			TimeNS: timer.GetElapsedNano(),
		})
		return
	}

	// Get location
	location, err := h.vehicleService.GetVehicleLocation(vehicleID)
	if err != nil || location == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:  "Location not found",
			Code:   domain.ErrCodeNotFound,
			TimeNS: timer.GetElapsedNano(),
		})
		return
	}

	c.JSON(http.StatusOK, models.VehicleLocationResponse{
		VehicleID:     vehicle.ID,
		VehicleNumber: vehicle.VehicleNumber,
		CurrentLocation: models.Location{
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
			Timestamp: location.Timestamp.Format(time.RFC3339),
		},
		CurrentGeofences: []models.CurrentGeofenceStatus{},
		TimeNS:           timer.GetElapsedNano(),
	})
}
