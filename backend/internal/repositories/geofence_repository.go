package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"geofence/internal/domain"
	"time"

	"github.com/lib/pq"
)

// PostgresGeofenceRepository implements GeofenceRepository
type PostgresGeofenceRepository struct {
	db *sql.DB
}

// NewPostgresGeofenceRepository creates a new geofence repository
func NewPostgresGeofenceRepository(db *sql.DB) domain.GeofenceRepository {
	return &PostgresGeofenceRepository{db: db}
}

// Create inserts a new geofence into the database
func (r *PostgresGeofenceRepository) Create(geofence *domain.Geofence) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert coordinates to flattened array for PostGIS
	flatCoords := pq.Float64Array{}
	for _, coord := range geofence.Coordinates {
		for _, val := range coord {
			flatCoords = append(flatCoords, val)
		}
	}

	query := `
		INSERT INTO geofences (id, name, description, coordinates, category, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		geofence.ID,
		geofence.Name,
		geofence.Description,
		flatCoords,
		string(geofence.Category),
		geofence.Status,
		now,
		now,
	).Scan(&geofence.ID, &geofence.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create geofence: %w", err)
	}

	geofence.UpdatedAt = now
	return nil
}

// GetByID retrieves a geofence by ID
func (r *PostgresGeofenceRepository) GetByID(id string) (*domain.Geofence, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, description, coordinates, category, status, created_at, updated_at
		FROM geofences
		WHERE id = $1
	`

	geofence := &domain.Geofence{}
	var flatCoords pq.Float64Array

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&geofence.ID,
		&geofence.Name,
		&geofence.Description,
		&flatCoords,
		&geofence.Category,
		&geofence.Status,
		&geofence.CreatedAt,
		&geofence.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrGeofenceNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get geofence: %w", err)
	}

	// Convert flattened array back to coordinates
	geofence.Coordinates = r.flattenedToCoordinates(flatCoords)

	return geofence, nil
}

// GetAll retrieves all geofences, optionally filtered by category
func (r *PostgresGeofenceRepository) GetAll(category *domain.GeofenceCategory) ([]*domain.Geofence, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, description, coordinates, category, status, created_at, updated_at
		FROM geofences
	`

	args := []interface{}{}

	if category != nil {
		query += `WHERE category = $1`
		args = append(args, string(*category))
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get geofences: %w", err)
	}
	defer rows.Close()

	var geofences []*domain.Geofence

	for rows.Next() {
		geofence := &domain.Geofence{}
		var flatCoords pq.Float64Array

		err := rows.Scan(
			&geofence.ID,
			&geofence.Name,
			&geofence.Description,
			&flatCoords,
			&geofence.Category,
			&geofence.Status,
			&geofence.CreatedAt,
			&geofence.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan geofence: %w", err)
		}

		geofence.Coordinates = r.flattenedToCoordinates(flatCoords)
		geofences = append(geofences, geofence)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating geofences: %w", err)
	}

	return geofences, nil
}

// Update updates an existing geofence
func (r *PostgresGeofenceRepository) Update(geofence *domain.Geofence) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	flatCoords := pq.Float64Array{}
	for _, coord := range geofence.Coordinates {
		for _, val := range coord {
			flatCoords = append(flatCoords, val)
		}
	}

	query := `
		UPDATE geofences
		SET name = $1, description = $2, coordinates = $3, category = $4, status = $5, updated_at = $6
		WHERE id = $7
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		geofence.Name,
		geofence.Description,
		flatCoords,
		string(geofence.Category),
		geofence.Status,
		now,
		geofence.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update geofence: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrGeofenceNotFound
	}

	geofence.UpdatedAt = now
	return nil
}

// Delete removes a geofence
func (r *PostgresGeofenceRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM geofences WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete geofence: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrGeofenceNotFound
	}

	return nil
}

// flattenedToCoordinates converts a flattened array to coordinates
func (r *PostgresGeofenceRepository) flattenedToCoordinates(flat pq.Float64Array) [][]float64 {
	if flat == nil || len(flat) == 0 {
		return [][]float64{}
	}

	coordinates := [][]float64{}
	for i := 0; i < len(flat); i += 2 {
		if i+1 < len(flat) {
			coordinates = append(coordinates, []float64{flat[i], flat[i+1]})
		}
	}
	return coordinates
}
