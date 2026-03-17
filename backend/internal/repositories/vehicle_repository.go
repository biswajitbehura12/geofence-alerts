package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"geofence/internal/domain"
	"time"
)

// PostgresVehicleRepository implements VehicleRepository
type PostgresVehicleRepository struct {
	db *sql.DB
}

// NewPostgresVehicleRepository creates a new vehicle repository
func NewPostgresVehicleRepository(db *sql.DB) domain.VehicleRepository {
	return &PostgresVehicleRepository{db: db}
}

// Create inserts a new vehicle
func (r *PostgresVehicleRepository) Create(vehicle *domain.Vehicle) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO vehicles (id, vehicle_number, driver_name, vehicle_type, phone, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		vehicle.ID,
		vehicle.VehicleNumber,
		vehicle.DriverName,
		vehicle.VehicleType,
		vehicle.Phone,
		vehicle.Status,
		now,
		now,
	).Scan(&vehicle.ID, &vehicle.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create vehicle: %w", err)
	}

	vehicle.UpdatedAt = now
	return nil
}

// GetByID retrieves a vehicle by ID
func (r *PostgresVehicleRepository) GetByID(id string) (*domain.Vehicle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, vehicle_number, driver_name, vehicle_type, phone, status, created_at, updated_at
		FROM vehicles
		WHERE id = $1
	`

	vehicle := &domain.Vehicle{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&vehicle.ID,
		&vehicle.VehicleNumber,
		&vehicle.DriverName,
		&vehicle.VehicleType,
		&vehicle.Phone,
		&vehicle.Status,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrVehicleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}

	return vehicle, nil
}

// GetAll retrieves all vehicles
func (r *PostgresVehicleRepository) GetAll() ([]*domain.Vehicle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, vehicle_number, driver_name, vehicle_type, phone, status, created_at, updated_at
		FROM vehicles
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []*domain.Vehicle

	for rows.Next() {
		vehicle := &domain.Vehicle{}
		err := rows.Scan(
			&vehicle.ID,
			&vehicle.VehicleNumber,
			&vehicle.DriverName,
			&vehicle.VehicleType,
			&vehicle.Phone,
			&vehicle.Status,
			&vehicle.CreatedAt,
			&vehicle.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}

		vehicles = append(vehicles, vehicle)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating vehicles: %w", err)
	}

	return vehicles, nil
}

// Update updates an existing vehicle
func (r *PostgresVehicleRepository) Update(vehicle *domain.Vehicle) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE vehicles
		SET vehicle_number = $1, driver_name = $2, vehicle_type = $3, phone = $4, status = $5, updated_at = $6
		WHERE id = $7
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		vehicle.VehicleNumber,
		vehicle.DriverName,
		vehicle.VehicleType,
		vehicle.Phone,
		vehicle.Status,
		now,
		vehicle.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrVehicleNotFound
	}

	vehicle.UpdatedAt = now
	return nil
}

// Delete removes a vehicle
func (r *PostgresVehicleRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM vehicles WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrVehicleNotFound
	}

	return nil
}

// GetByVehicleNumber retrieves a vehicle by vehicle number
func (r *PostgresVehicleRepository) GetByVehicleNumber(vehicleNumber string) (*domain.Vehicle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, vehicle_number, driver_name, vehicle_type, phone, status, created_at, updated_at
		FROM vehicles
		WHERE vehicle_number = $1
	`

	vehicle := &domain.Vehicle{}
	err := r.db.QueryRowContext(ctx, query, vehicleNumber).Scan(
		&vehicle.ID,
		&vehicle.VehicleNumber,
		&vehicle.DriverName,
		&vehicle.VehicleType,
		&vehicle.Phone,
		&vehicle.Status,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrVehicleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}

	return vehicle, nil
}

// PostgresVehicleLocationRepository implements VehicleLocationRepository
type PostgresVehicleLocationRepository struct {
	db *sql.DB
}

// NewPostgresVehicleLocationRepository creates a new vehicle location repository
func NewPostgresVehicleLocationRepository(db *sql.DB) domain.VehicleLocationRepository {
	return &PostgresVehicleLocationRepository{db: db}
}

// SaveLocation saves a vehicle location
func (r *PostgresVehicleLocationRepository) SaveLocation(location *domain.VehicleLocation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO vehicle_locations (id, vehicle_id, latitude, longitude, timestamp, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		location.ID,
		location.VehicleID,
		location.Latitude,
		location.Longitude,
		location.Timestamp,
		now,
	).Scan(&location.ID, &location.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to save location: %w", err)
	}

	return nil
}

// GetLatestLocation retrieves the latest location for a vehicle
func (r *PostgresVehicleLocationRepository) GetLatestLocation(vehicleID string) (*domain.VehicleLocation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, vehicle_id, latitude, longitude, timestamp, created_at
		FROM vehicle_locations
		WHERE vehicle_id = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`

	location := &domain.VehicleLocation{}
	err := r.db.QueryRowContext(ctx, query, vehicleID).Scan(
		&location.ID,
		&location.VehicleID,
		&location.Latitude,
		&location.Longitude,
		&location.Timestamp,
		&location.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}

	return location, nil
}

// GetLocationHistory retrieves location history for a vehicle
func (r *PostgresVehicleLocationRepository) GetLocationHistory(vehicleID string, limit int, offset int) ([]*domain.VehicleLocation, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get total count
	countQuery := `SELECT COUNT(*) FROM vehicle_locations WHERE vehicle_id = $1`
	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, vehicleID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get count: %w", err)
	}

	query := `
		SELECT id, vehicle_id, latitude, longitude, timestamp, created_at
		FROM vehicle_locations
		WHERE vehicle_id = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, vehicleID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get location history: %w", err)
	}
	defer rows.Close()

	var locations []*domain.VehicleLocation

	for rows.Next() {
		location := &domain.VehicleLocation{}
		err := rows.Scan(
			&location.ID,
			&location.VehicleID,
			&location.Latitude,
			&location.Longitude,
			&location.Timestamp,
			&location.CreatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan location: %w", err)
		}

		locations = append(locations, location)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating locations: %w", err)
	}

	return locations, totalCount, nil
}
