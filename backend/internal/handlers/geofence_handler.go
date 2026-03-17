package handlers

import (
	"geofence/internal/domain"
	"geofence/internal/models"
	"geofence/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GeofenceHandler handles geofence-related HTTP requests
type GeofenceHandler struct {
	geofenceService domain.GeofenceService
}

// NewGeofenceHandler creates a new geofence handler
func NewGeofenceHandler(geofenceService domain.GeofenceService) *GeofenceHandler {
	return &GeofenceHandler{
		geofenceService: geofenceService,
	}
}

// CreateGeofence handles POST /geofences
func (h *GeofenceHandler) CreateGeofence(c *gin.Context) {
	timer := services.NewTimeHelper()

	var req models.CreateGeofenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:  "Invalid request body",
			Code:   domain.ErrCodeInvalidInput,
			TimeNS: timer.GetElapsedNano(),
		})
		return
	}

	// Create domain geofence
	geofence := &domain.Geofence{
		Name:        req.Name,
		Description: req.Description,
		Coordinates: req.Coordinates,
		Category:    domain.GeofenceCategory(req.Category),
		Status:      "active",
	}

	// Create geofence
	result, err := h.geofenceService.CreateGeofence(geofence)
	if err != nil {
		code := domain.ErrCodeInternalError
		status := http.StatusInternalServerError

		if domErr, ok := err.(*domain.DomainError); ok {
			code = domErr.Code
			if domErr.Code == domain.ErrCodeInvalidCoords {
				status = http.StatusBadRequest
			}
		}

		c.JSON(status, models.ErrorResponse{
			Error:  err.Error(),
			Code:   code,
			TimeNS: timer.GetElapsedNano(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.GeofenceResponse{
		ID:     result.ID,
		Name:   result.Name,
		Status: result.Status,
		TimeNS: timer.GetElapsedNano(),
	})
}

// GetGeofences handles GET /geofences
func (h *GeofenceHandler) GetGeofences(c *gin.Context) {
	timer := services.NewTimeHelper()

	// Get category filter
	categoryStr := c.Query("category")
	var category *domain.GeofenceCategory
	if categoryStr != "" {
		cat := domain.GeofenceCategory(categoryStr)
		category = &cat
	}

	// Get geofences
	geofences, err := h.geofenceService.GetGeofences(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:  err.Error(),
			Code:   domain.ErrCodeInternalError,
			TimeNS: timer.GetElapsedNano(),
		})
		return
	}

	// Convert to response
	details := []*models.GeofenceDetail{}
	for _, g := range geofences {
		details = append(details, &models.GeofenceDetail{
			ID:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			Coordinates: g.Coordinates,
			Category:    string(g.Category),
			CreatedAt:   g.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, models.GeofenceListResponse{
		Geofences: details,
		TimeNS:    timer.GetElapsedNano(),
	})
}
