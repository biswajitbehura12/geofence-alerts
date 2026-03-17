package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"geofence/internal/domain"
	"time"
)

// PostgresAlertRepository implements AlertRepository
type PostgresAlertRepository struct {
	db *sql.DB
}

// NewPostgresAlertRepository creates a new alert repository
func NewPostgresAlertRepository(db *sql.DB) domain.AlertRepository {
	return &PostgresAlertRepository{db: db}
}

// CreateRule creates a new alert rule
func (r *PostgresAlertRepository) CreateRule(rule *domain.AlertRule) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO alert_rules (id, geofence_id, vehicle_id, event_type, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		rule.ID,
		rule.GeofenceID,
		rule.VehicleID,
		string(rule.EventType),
		rule.Status,
		now,
		now,
	).Scan(&rule.ID, &rule.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create alert rule: %w", err)
	}

	rule.UpdatedAt = now
	return nil
}

// GetRuleByID retrieves an alert rule by ID
func (r *PostgresAlertRepository) GetRuleByID(id string) (*domain.AlertRule, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, geofence_id, vehicle_id, event_type, status, created_at, updated_at
		FROM alert_rules
		WHERE id = $1
	`

	rule := &domain.AlertRule{}
	var eventType string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rule.ID,
		&rule.GeofenceID,
		&rule.VehicleID,
		&eventType,
		&rule.Status,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrAlertNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get alert rule: %w", err)
	}

	rule.EventType = domain.EventType(eventType)
	return rule, nil
}

// GetAllRules retrieves all alert rules, optionally filtered
func (r *PostgresAlertRepository) GetAllRules(geofenceID, vehicleID *string) ([]*domain.AlertRule, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, geofence_id, vehicle_id, event_type, status, created_at, updated_at
		FROM alert_rules
		WHERE 1=1
	`

	args := []interface{}{}
	argNum := 1

	if geofenceID != nil {
		query += fmt.Sprintf(` AND geofence_id = $%d`, argNum)
		args = append(args, *geofenceID)
		argNum++
	}

	if vehicleID != nil {
		query += fmt.Sprintf(` AND vehicle_id = $%d`, argNum)
		args = append(args, *vehicleID)
		argNum++
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert rules: %w", err)
	}
	defer rows.Close()

	var rules []*domain.AlertRule

	for rows.Next() {
		rule := &domain.AlertRule{}
		var eventType string

		err := rows.Scan(
			&rule.ID,
			&rule.GeofenceID,
			&rule.VehicleID,
			&eventType,
			&rule.Status,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan alert rule: %w", err)
		}

		rule.EventType = domain.EventType(eventType)
		rules = append(rules, rule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating alert rules: %w", err)
	}

	return rules, nil
}

// DeleteRule deletes an alert rule
func (r *PostgresAlertRepository) DeleteRule(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM alert_rules WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete alert rule: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrAlertNotFound
	}

	return nil
}

// SaveEvent saves a geofence event
func (r *PostgresAlertRepository) SaveEvent(event *domain.GeofenceEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO geofence_events (id, vehicle_id, geofence_id, event_type, latitude, longitude, timestamp, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		event.ID,
		event.VehicleID,
		event.GeofenceID,
		string(event.EventType),
		event.Latitude,
		event.Longitude,
		event.Timestamp,
		now,
	).Scan(&event.ID, &event.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to save event: %w", err)
	}

	return nil
}

// GetEventHistory retrieves geofence event history
func (r *PostgresAlertRepository) GetEventHistory(vehicleID, geofenceID *string, startDate, endDate *time.Time, limit, offset int) ([]*domain.GeofenceEvent, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Build count query
	countQuery := `SELECT COUNT(*) FROM geofence_events WHERE 1=1`
	queryBase := `SELECT id, vehicle_id, geofence_id, event_type, latitude, longitude, timestamp, created_at FROM geofence_events WHERE 1=1`

	args := []interface{}{}
	argNum := 1

	if vehicleID != nil {
		condition := fmt.Sprintf(` AND vehicle_id = $%d`, argNum)
		countQuery += condition
		queryBase += condition
		args = append(args, *vehicleID)
		argNum++
	}

	if geofenceID != nil {
		condition := fmt.Sprintf(` AND geofence_id = $%d`, argNum)
		countQuery += condition
		queryBase += condition
		args = append(args, *geofenceID)
		argNum++
	}

	if startDate != nil {
		condition := fmt.Sprintf(` AND timestamp >= $%d`, argNum)
		countQuery += condition
		queryBase += condition
		args = append(args, *startDate)
		argNum++
	}

	if endDate != nil {
		condition := fmt.Sprintf(` AND timestamp <= $%d`, argNum)
		countQuery += condition
		queryBase += condition
		args = append(args, *endDate)
		argNum++
	}

	// Get total count
	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get count: %w", err)
	}

	// Get paginated results
	queryBase += ` ORDER BY timestamp DESC LIMIT $` + fmt.Sprintf("%d OFFSET $%d", argNum, argNum+1)
	listArgs := append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, queryBase, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get event history: %w", err)
	}
	defer rows.Close()

	var events []*domain.GeofenceEvent

	for rows.Next() {
		event := &domain.GeofenceEvent{}
		var eventType string

		err := rows.Scan(
			&event.ID,
			&event.VehicleID,
			&event.GeofenceID,
			&eventType,
			&event.Latitude,
			&event.Longitude,
			&event.Timestamp,
			&event.CreatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan event: %w", err)
		}

		event.EventType = domain.EventType(eventType)
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating events: %w", err)
	}

	return events, totalCount, nil
}
