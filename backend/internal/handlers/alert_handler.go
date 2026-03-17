package handlers

import (
	"geofence/internal/domain"
	"geofence/internal/models"
	"geofence/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AlertHandler handles alert-related HTTP requests
type AlertHandler struct {
	alertService domain.AlertService
}

// NewAlertHandler creates a new alert handler
func NewAlertHandler(alertService domain.AlertService) *AlertHandler {
	return &AlertHandler{
		alertService: alertService,
	}
}

// ConfigureAlert handles POST /alerts/configure
func (h *AlertHandler) ConfigureAlert(c *gin.Context) {
	timer := services.NewTimeHelper()

	var req models.ConfigureAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:  "Invalid request body",
			Code:   domain.ErrCodeInvalidInput,
			TimeNS: timer.GetElapsedNano(),
		})
		return
	}

	// Configure alert
	rule, err := h.alertService.ConfigureAlert(req.GeofenceID, req.VehicleID, domain.EventType(req.EventType))
	if err != nil {
		code := domain.ErrCodeInternalError
		status := http.StatusInternalServerError

		if err == domain.ErrGeofenceNotFound || err == domain.ErrVehicleNotFound {
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

	c.JSON(http.StatusCreated, models.AlertResponse{
		AlertID:    rule.ID,
		GeofenceID: rule.GeofenceID,
		VehicleID:  rule.VehicleID,
		EventType:  string(rule.EventType),
		Status:     rule.Status,
		TimeNS:     timer.GetElapsedNano(),
	})
}

// GetAlerts handles GET /alerts
func (h *AlertHandler) GetAlerts(c *gin.Context) {
	timer := services.NewTimeHelper()

	// Get filters
	geofenceID := c.Query("geofence_id")
	vehicleID := c.Query("vehicle_id")

	var geoFilter *string
	var vehFilter *string

	if geofenceID != "" {
		geoFilter = &geofenceID
	}

	if vehicleID != "" {
		vehFilter = &vehicleID
	}

	// Get alerts
	rules, err := h.alertService.GetAlerts(geoFilter, vehFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:  err.Error(),
			Code:   domain.ErrCodeInternalError,
			TimeNS: timer.GetElapsedNano(),
		})
		return
	}

	// TODO: Get geofence and vehicle names for each rule
	// For now, return basic info
	details := []*models.AlertDetail{}
	for _, r := range rules {
		details = append(details, &models.AlertDetail{
			AlertID:       r.ID,
			GeofenceID:    r.GeofenceID,
			GeofenceName:  r.GeofenceID, // TODO: fetch actual name
			VehicleID:     r.VehicleID,
			VehicleNumber: r.VehicleID, // TODO: fetch actual number
			EventType:     string(r.EventType),
			Status:        r.Status,
			CreatedAt:     r.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, models.AlertListResponse{
		Alerts: details,
		TimeNS: timer.GetElapsedNano(),
	})
}

// GetViolationHistory handles GET /violations/history
func (h *AlertHandler) GetViolationHistory(c *gin.Context) {
	timer := services.NewTimeHelper()

	// Get filters
	vehicleID := c.Query("vehicle_id")
	geofenceID := c.Query("geofence_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	var vehFilter *string
	var geoFilter *string
	var startDate *time.Time
	var endDate *time.Time

	if vehicleID != "" {
		vehFilter = &vehicleID
	}

	if geofenceID != "" {
		geoFilter = &geofenceID
	}

	if startDateStr != "" {
		t, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:  "Invalid start_date format",
				Code:   domain.ErrCodeInvalidInput,
				TimeNS: timer.GetElapsedNano(),
			})
			return
		}
		startDate = &t
	}

	if endDateStr != "" {
		t, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:  "Invalid end_date format",
				Code:   domain.ErrCodeInvalidInput,
				TimeNS: timer.GetElapsedNano(),
			})
			return
		}
		endDate = &t
	}

	// Parse limit and offset
	limit := 50
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	// Get violation history
	events, total, err := h.alertService.GetViolationHistory(vehFilter, geoFilter, startDate, endDate, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:  err.Error(),
			Code:   domain.ErrCodeInternalError,
			TimeNS: timer.GetElapsedNano(),
		})
		return
	}

	// Convert to response
	details := []*models.ViolationDetail{}
	for _, e := range events {
		details = append(details, &models.ViolationDetail{
			ID:            e.ID,
			VehicleID:     e.VehicleID,
			VehicleNumber: e.VehicleID, // TODO: fetch actual vehicle number
			GeofenceID:    e.GeofenceID,
			GeofenceName:  e.GeofenceID, // TODO: fetch actual geofence name
			EventType:     string(e.EventType),
			Latitude:      e.Latitude,
			Longitude:     e.Longitude,
			Timestamp:     e.Timestamp.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, models.ViolationHistoryResponse{
		Violations: details,
		TotalCount: total,
		TimeNS:     timer.GetElapsedNano(),
	})
}
