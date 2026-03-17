package main

import (
	"database/sql"
	"fmt"
	"geofence/config"
	"geofence/internal/handlers"
	"geofence/internal/repositories"
	"geofence/internal/services"
	"log"

	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := initializeDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Check database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	geofenceRepo := repositories.NewPostgresGeofenceRepository(db)
	vehicleRepo := repositories.NewPostgresVehicleRepository(db)
	vehicleLocRepo := repositories.NewPostgresVehicleLocationRepository(db)
	alertRepo := repositories.NewPostgresAlertRepository(db)

	// Initialize alert publisher
	alertPublisher := services.NewInMemoryAlertPublisher()

	// Initialize services
	geofenceService := services.NewGeofenceService(geofenceRepo)
	vehicleService := services.NewVehicleService(vehicleRepo, vehicleLocRepo, geofenceRepo, geofenceService)
	alertService := services.NewAlertService(alertRepo, vehicleService, geofenceService, geofenceRepo, alertPublisher)

	// Initialize handlers
	geofenceHandler := handlers.NewGeofenceHandler(geofenceService)
	vehicleHandler := handlers.NewVehicleHandler(vehicleService, alertService, geofenceService)
	alertHandler := handlers.NewAlertHandler(alertService)
	wsHandler := handlers.NewWebSocketHandler(alertPublisher)

	// Setup Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Geofence endpoints
	router.POST("/geofences", geofenceHandler.CreateGeofence)
	router.GET("/geofences", geofenceHandler.GetGeofences)

	// Vehicle endpoints
	router.POST("/vehicles", vehicleHandler.RegisterVehicle)
	router.GET("/vehicles", vehicleHandler.GetVehicles)
	router.POST("/vehicles/location", vehicleHandler.UpdateLocation)
	router.GET("/vehicles/location/:vehicle_id", vehicleHandler.GetLocation)

	// Alert endpoints
	router.POST("/alerts/configure", alertHandler.ConfigureAlert)
	router.GET("/alerts", alertHandler.GetAlerts)
	router.GET("/violations/history", alertHandler.GetViolationHistory)

	// WebSocket endpoint
	router.GET("/ws/alerts", func(c *gin.Context) {
		wsHandler.HandleWebSocket(c)
	})

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initializeDatabase initializes the database connection
func initializeDatabase(cfg config.DatabaseConfig) (*sql.DB, error) {
	connStr := cfg.GetConnectionURL()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}

// runMigrations runs database migrations
func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS geofences (
			id VARCHAR(50) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			coordinates FLOAT8[] NOT NULL,
			category VARCHAR(50) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS vehicles (
			id VARCHAR(50) PRIMARY KEY,
			vehicle_number VARCHAR(50) NOT NULL UNIQUE,
			driver_name VARCHAR(255) NOT NULL,
			vehicle_type VARCHAR(50) NOT NULL,
			phone VARCHAR(20) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS vehicle_locations (
			id VARCHAR(50) PRIMARY KEY,
			vehicle_id VARCHAR(50) NOT NULL REFERENCES vehicles(id),
			latitude FLOAT8 NOT NULL,
			longitude FLOAT8 NOT NULL,
			timestamp TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS alert_rules (
			id VARCHAR(50) PRIMARY KEY,
			geofence_id VARCHAR(50) NOT NULL REFERENCES geofences(id),
			vehicle_id VARCHAR(50),
			event_type VARCHAR(50) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS geofence_events (
			id VARCHAR(50) PRIMARY KEY,
			vehicle_id VARCHAR(50) NOT NULL REFERENCES vehicles(id),
			geofence_id VARCHAR(50) NOT NULL REFERENCES geofences(id),
			event_type VARCHAR(50) NOT NULL,
			latitude FLOAT8 NOT NULL,
			longitude FLOAT8 NOT NULL,
			timestamp TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE INDEX IF NOT EXISTS idx_vehicles_id ON vehicle_locations(vehicle_id)`,
		`CREATE INDEX IF NOT EXISTS idx_geofence_events_vehicle ON geofence_events(vehicle_id)`,
		`CREATE INDEX IF NOT EXISTS idx_geofence_events_geofence ON geofence_events(geofence_id)`,
		`CREATE INDEX IF NOT EXISTS idx_alert_rules_geofence ON alert_rules(geofence_id)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	log.Println("Migrations completed successfully")
	return nil
}
