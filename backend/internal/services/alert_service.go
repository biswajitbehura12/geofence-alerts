package services

import (
	"fmt"
	"geofence/internal/domain"
	"github.com/google/uuid"
	"sync"
	"time"
)

// AlertService implements domain.AlertService
type AlertService struct {
	alertRepo       domain.AlertRepository
	vehicleService  domain.VehicleService
	geofenceService domain.GeofenceService
	geofenceRepo    domain.GeofenceRepository
	alertPublisher  domain.AlertPublisher
}

// NewAlertService creates a new alert service
func NewAlertService(
	alertRepo domain.AlertRepository,
	vehicleService domain.VehicleService,
	geofenceService domain.GeofenceService,
	geofenceRepo domain.GeofenceRepository,
	alertPublisher domain.AlertPublisher,
) domain.AlertService {
	return &AlertService{
		alertRepo:       alertRepo,
		vehicleService:  vehicleService,
		geofenceService: geofenceService,
		geofenceRepo:    geofenceRepo,
		alertPublisher:  alertPublisher,
	}
}

// ConfigureAlert configures a new alert rule
func (s *AlertService) ConfigureAlert(geofenceID, vehicleID string, eventType domain.EventType) (*domain.AlertRule, error) {
	// Validate geofence exists
	_, err := s.geofenceRepo.GetByID(geofenceID)
	if err != nil {
		return nil, err
	}

	// If vehicleID is provided, validate vehicle exists
	if vehicleID != "" {
		_, err = s.vehicleService.GetVehicleByID(vehicleID)
		if err != nil {
			return nil, err
		}
	}

	// Validate event type
	if eventType != domain.EventEntry && eventType != domain.EventExit && eventType != domain.EventBoth {
		return nil, fmt.Errorf("invalid event type: %s", eventType)
	}

	// Create alert rule
	rule := &domain.AlertRule{
		ID:         "alert_" + uuid.New().String()[:12],
		GeofenceID: geofenceID,
		VehicleID:  vehicleID,
		EventType:  eventType,
		Status:     "active",
	}

	if err := s.alertRepo.CreateRule(rule); err != nil {
		return nil, fmt.Errorf("failed to configure alert: %w", err)
	}

	return rule, nil
}

// GetAlerts retrieves alert rules
func (s *AlertService) GetAlerts(geofenceID, vehicleID *string) ([]*domain.AlertRule, error) {
	return s.alertRepo.GetAllRules(geofenceID, vehicleID)
}

// GetViolationHistory retrieves violation history
func (s *AlertService) GetViolationHistory(vehicleID, geofenceID *string, startDate, endDate *time.Time, limit, offset int) ([]*domain.GeofenceEvent, int64, error) {
	// Default limit
	if limit == 0 {
		limit = 50
	}

	// Max limit
	if limit > 500 {
		limit = 500
	}

	return s.alertRepo.GetEventHistory(vehicleID, geofenceID, startDate, endDate, limit, offset)
}

// CheckAndTriggerAlerts checks if alerts should be triggered and publishes them
func (s *AlertService) CheckAndTriggerAlerts(vehicle *domain.Vehicle, currentGeofences []*domain.Geofence, eventType domain.EventType) ([]*domain.RealTimeAlert, error) {
	var alerts []*domain.RealTimeAlert

	// For each triggered geofence
	for _, geofence := range currentGeofences {
		// Get matching alert rules
		rules, err := s.alertRepo.GetAllRules(&geofence.ID, nil)
		if err != nil {
			continue
		}

		// Check each rule
		for _, rule := range rules {
			// Check if alert applies to this vehicle (empty = all vehicles)
			if rule.VehicleID != "" && rule.VehicleID != vehicle.ID {
				continue
			}

			// Check if event type matches
			if rule.EventType != eventType && rule.EventType != domain.EventBoth {
				continue
			}

			// Get current location
			location, _ := s.vehicleService.GetVehicleLocation(vehicle.ID)
			if location == nil {
				continue
			}

			// Create real-time alert
			alert := &domain.RealTimeAlert{
				EventID:   "evt_" + uuid.New().String()[:12],
				EventType: eventType,
				Timestamp: time.Now(),
				Vehicle: domain.VehicleAlertInfo{
					VehicleID:    vehicle.ID,
					VehicleNumber: vehicle.VehicleNumber,
					DriverName:   vehicle.DriverName,
				},
				Geofence: domain.GeofenceAlertInfo{
					GeofenceID:   geofence.ID,
					GeofenceName: geofence.Name,
					Category:    geofence.Category,
				},
				Location: domain.LocationAlertInfo{
					Latitude:  location.Latitude,
					Longitude: location.Longitude,
				},
			}

			alerts = append(alerts, alert)

			// Save event asynchronously
			go func(a *domain.RealTimeAlert) {
				event := &domain.GeofenceEvent{
					ID:         uuid.New().String(),
					VehicleID:  a.Vehicle.VehicleID,
					GeofenceID: a.Geofence.GeofenceID,
					EventType:  a.EventType,
					Latitude:   a.Location.Latitude,
					Longitude:  a.Location.Longitude,
					Timestamp:  a.Timestamp,
				}
				_ = s.alertRepo.SaveEvent(event)
			}(alert)

			// Publish alert asynchronously
			go func(a *domain.RealTimeAlert) {
				_ = s.alertPublisher.PublishAlert(a)
			}(alert)
		}
	}

	return alerts, nil
}

// InMemoryAlertPublisher implements AlertPublisher using channels
type InMemoryAlertPublisher struct {
	subscribers map[string]chan *domain.RealTimeAlert
	mu          sync.RWMutex
}

// NewInMemoryAlertPublisher creates a new in-memory alert publisher
func NewInMemoryAlertPublisher() domain.AlertPublisher {
	return &InMemoryAlertPublisher{
		subscribers: make(map[string]chan *domain.RealTimeAlert),
	}
}

// PublishAlert publishes an alert to all subscribers
func (p *InMemoryAlertPublisher) PublishAlert(alert *domain.RealTimeAlert) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, ch := range p.subscribers {
		select {
		case ch <- alert:
		default:
			// Channel full, skip
		}
	}

	return nil
}

// Subscribe subscribes to alerts
func (p *InMemoryAlertPublisher) Subscribe(ch chan *domain.RealTimeAlert) string {
	p.mu.Lock()
	defer p.mu.Unlock()

	subscriptionID := uuid.New().String()
	p.subscribers[subscriptionID] = ch

	return subscriptionID
}

// Unsubscribe unsubscribes from alerts
func (p *InMemoryAlertPublisher) Unsubscribe(subscriptionID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if ch, exists := p.subscribers[subscriptionID]; exists {
		close(ch)
		delete(p.subscribers, subscriptionID)
	}
}
