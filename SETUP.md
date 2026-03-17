# Geofence & Vehicle Tracking System - Setup Guide

Complete documentation for setting up and running the geofencing and real-time vehicle tracking system.

## 📋 Table of Contents

1. [Prerequisites](#prerequisites)
2. [Local Setup](#local-setup)
3. [Docker Setup](#docker-setup)
4. [Database Configuration](#database-configuration)
5. [Running the Application](#running-the-application)
6. [API Documentation](#api-documentation)
7. [Frontend Usage Guide](#frontend-usage-guide)
8. [Architecture Overview](#architecture-overview)
9. [Testing](#testing)
10. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Software

- **Go 1.21+** - [Download](https://golang.org/dl/)
- **Node.js 18+** - [Download](https://nodejs.org/)
- **PostgreSQL 13+** - [Download](https://www.postgresql.org/download/)
- **Docker & Docker Compose** (optional but recommended)
- **Git** - [Download](https://git-scm.com/)

### Recommended Tools

- **Postman** - API testing tool
- **VS Code** - Code editor
- **DBeaver** - Database GUI client

---

## Local Setup

### 1. Clone Repository

```bash
git clone <repository-url>
cd geofence
```

### 2. Setup Backend

```bash
cd backend

# Install Go dependencies
go mod download

# Copy environment configuration
cp .env.example .env

# Edit .env with your database credentials
nano .env
```

### 3. Setup Frontend

```bash
cd ../frontend

# Install npm dependencies
npm install

# Create environment file
cp .env.example .env.local

# Edit environment variables
nano .env.local
```

### 4. Setup Database

#### Option A: Using Docker (Recommended)

```bash
# Start PostgreSQL container
docker run -d \
  --name geofence-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=geofence \
  -p 5432:5432 \
  postgres:15-alpine
```

#### Option B: Local PostgreSQL Installation

```bash
# Create database
createdb -U postgres geofence

# The application will automatically create tables on startup
```

### 5. Run Backend

```bash
cd backend

# Set environment variables (or use .env)
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=geofence

# Run the application
go run ./cmd/main.go
```

Backend will start at http://localhost:8080

### 6. Run Frontend

```bash
cd frontend

# Start development server
npm start
```

Frontend will open at http://localhost:3000

---

## Docker Setup

### Using Docker Compose (Easiest)

```bash
# Copy environment file
cp .env.example .env

# Edit environment variables if needed
nano .env

# Build and start all services
docker-compose up --build

# To run in background
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Manual Docker Build

#### Build Backend Image

```bash
cd backend

# Build image
docker build -f docker/Dockerfile -t geofence-backend:latest .

# Run container
docker run -d \
  --name geofence-backend \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=geofence \
  -p 8080:8080 \
  geofence-backend:latest
```

#### Build Frontend Image

```bash
cd frontend

# Build image
docker build -t geofence-frontend:latest .

# Run container
docker run -d \
  --name geofence-frontend \
  -e REACT_APP_API_URL=http://localhost:8080 \
  -p 3000:3000 \
  geofence-frontend:latest
```

---

## Database Configuration

### Schema Overview

The application automatically creates the following tables:

- **geofences** - Virtual boundaries
- **vehicles** - Registered vehicles
- **vehicle_locations** - Location history
- **alert_rules** - Alert configuration
- **geofence_events** - Violation history

### Manual Database Setup

If tables aren't created automatically:

```sql
-- Connect to database
psql -U postgres -d geofence

-- Tables are created automatically by the application
-- No manual SQL setup required
```

### Connection String Format

```
postgresql://user:password@host:port/database?sslmode=disable
```

Example:

```
postgresql://postgres:postgres@localhost:5432/geofence?sslmode=disable
```

---

## Running the Application

### Production Deployment

#### Backend on Cloud (e.g., AWS EC2)

```bash
# SSH into instance
ssh -i key.pem ec2-user@your-instance-ip

# Install Go
wget https://golang.org/dl/go1.21.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.linux-amd64.tar.gz

# Clone and setup
git clone <repo-url>
cd geofence/backend

# Build binary
go build -o geofence ./cmd/main.go

# Run with proper environment
export DB_HOST=your-rds-endpoint
export DB_USER=postgres
export DB_PASSWORD=your-password
export DB_NAME=geofence

./geofence
```

#### Frontend on Vercel

```bash
# Install Vercel CLI
npm install -g vercel

# Login to Vercel
vercel login

# Deploy
cd frontend
vercel

# Set environment variables in Vercel dashboard
# REACT_APP_API_URL=your-backend-url
# REACT_APP_WS_URL=your-websocket-url
```

---

## API Documentation

### Base URL

```
http://localhost:8080
```

### Authentication

Currently no authentication required. Add JWT/API keys for production.

### Response Format

All responses include execution time in nanoseconds:

```json
{
  "data": {...},
  "time_ns": "1234567"
}
```

### 1. Geofence Endpoints

#### Create Geofence

```bash
POST /geofences

Request:
{
  "name": "Downtown Zone",
  "description": "Main delivery area",
  "coordinates": [
    [37.7749, -122.4194],
    [37.7849, -122.4194],
    [37.7849, -122.4094],
    [37.7749, -122.4094],
    [37.7749, -122.4194]
  ],
  "category": "delivery_zone"
}

Response:
{
  "id": "geo_abc123",
  "name": "Downtown Zone",
  "status": "active",
  "time_ns": "1234567"
}
```

#### Get All Geofences

```bash
GET /geofences?category=delivery_zone

Response:
{
  "geofences": [
    {
      "id": "geo_123",
      "name": "Downtown Zone",
      "description": "Main delivery area",
      "coordinates": [[37.7749, -122.4194], ...],
      "category": "delivery_zone",
      "created_at": "2025-01-15T10:30:00Z"
    }
  ],
  "time_ns": "987654"
}
```

### 2. Vehicle Endpoints

#### Register Vehicle

```bash
POST /vehicles

Request:
{
  "vehicle_number": "KA-01-AB-1234",
  "driver_name": "John Doe",
  "vehicle_type": "truck",
  "phone": "+1234567890"
}

Response:
{
  "id": "veh_456",
  "vehicle_number": "KA-01-AB-1234",
  "status": "active",
  "time_ns": "1123456"
}
```

#### Get All Vehicles

```bash
GET /vehicles

Response:
{
  "vehicles": [
    {
      "id": "veh_456",
      "vehicle_number": "KA-01-AB-1234",
      "driver_name": "John Doe",
      "vehicle_type": "truck",
      "phone": "+1234567890",
      "status": "active",
      "created_at": "2025-01-15T09:00:00Z"
    }
  ],
  "time_ns": "876543"
}
```

#### Update Vehicle Location

```bash
POST /vehicles/location

Request:
{
  "vehicle_id": "veh_456",
  "latitude": 37.7849,
  "longitude": -122.4194,
  "timestamp": "2025-01-15T10:35:00Z"
}

Response:
{
  "vehicle_id": "veh_456",
  "location_updated": true,
  "current_geofences": [
    {
      "geofence_id": "geo_123",
      "geofence_name": "Downtown Zone",
      "status": "inside"
    }
  ],
  "time_ns": "2345678"
}
```

#### Get Vehicle Location

```bash
GET /vehicles/location/{vehicle_id}

Response:
{
  "vehicle_id": "veh_456",
  "vehicle_number": "KA-01-AB-1234",
  "current_location": {
    "latitude": 37.7849,
    "longitude": -122.4194,
    "timestamp": "2025-01-15T10:35:00Z"
  },
  "current_geofences": [
    {
      "geofence_id": "geo_123",
      "geofence_name": "Downtown Zone",
      "category": "delivery_zone"
    }
  ],
  "time_ns": "876543"
}
```

### 3. Alert Endpoints

#### Configure Alert

```bash
POST /alerts/configure

Request:
{
  "geofence_id": "geo_123",
  "vehicle_id": "veh_456",
  "event_type": "entry"
}

Response:
{
  "alert_id": "alert_789",
  "geofence_id": "geo_123",
  "vehicle_id": "veh_456",
  "event_type": "entry",
  "status": "active",
  "time_ns": "1567890"
}
```

#### Get Alerts

```bash
GET /alerts?geofence_id=geo_123&vehicle_id=veh_456

Response:
{
  "alerts": [
    {
      "alert_id": "alert_789",
      "geofence_id": "geo_123",
      "geofence_name": "Downtown Zone",
      "vehicle_id": "veh_456",
      "vehicle_number": "KA-01-AB-1234",
      "event_type": "entry",
      "status": "active",
      "created_at": "2025-01-15T09:15:00Z"
    }
  ],
  "time_ns": "654321"
}
```

#### Get Violation History

```bash
GET /violations/history?vehicle_id=veh_456&limit=100&offset=0&start_date=2025-01-01T00:00:00Z

Response:
{
  "violations": [
    {
      "id": "viol_111",
      "vehicle_id": "veh_456",
      "vehicle_number": "KA-01-AB-1234",
      "geofence_id": "geo_123",
      "geofence_name": "Downtown Zone",
      "event_type": "entry",
      "latitude": 37.7849,
      "longitude": -122.4194,
      "timestamp": "2025-01-15T10:35:00Z"
    }
  ],
  "total_count": 245,
  "time_ns": "3456789"
}
```

### 4. WebSocket Endpoint

#### Connect to Real-time Alerts

```bash
WebSocket: ws://localhost:8080/ws/alerts

# Message Format (received):
{
  "event_id": "evt_999",
  "event_type": "entry",
  "timestamp": "2025-01-15T10:35:00Z",
  "vehicle": {
    "vehicle_id": "veh_456",
    "vehicle_number": "KA-01-AB-1234",
    "driver_name": "John Doe"
  },
  "geofence": {
    "geofence_id": "geo_123",
    "geofence_name": "Downtown Zone",
    "category": "delivery_zone"
  },
  "location": {
    "latitude": 37.7849,
    "longitude": -122.4194
  }
}
```

### Testing with Curl

```bash
# Create geofence
curl -X POST http://localhost:8080/geofences \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Zone",
    "description": "Test",
    "coordinates": [[37.7, -122.4], [37.8, -122.4], [37.8, -122.3], [37.7, -122.3], [37.7, -122.4]],
    "category": "delivery_zone"
  }'

# Register vehicle
curl -X POST http://localhost:8080/vehicles \
  -H "Content-Type: application/json" \
  -d '{
    "vehicle_number": "TEST-001",
    "driver_name": "Test Driver",
    "vehicle_type": "truck",
    "phone": "+1234567890"
  }'

# Get all vehicles
curl http://localhost:8080/vehicles

# Update location
curl -X POST http://localhost:8080/vehicles/location \
  -H "Content-Type: application/json" \
  -d '{
    "vehicle_id": "veh_xxx",
    "latitude": 37.75,
    "longitude": -122.4,
    "timestamp": "2025-01-15T10:35:00Z"
  }'
```

---

## Frontend Usage Guide

### Dashboard

- View system overview with statistics
- See real-time geofence and vehicle information
- Monitor active WebSocket connection status
- Quick access to recent alerts

### Geofence Management

- **Create** new geofences with polygonal boundaries
- **View** all geofences on interactive map
- **Filter** by category (delivery zone, restricted zone, toll zone, customer area)
- **Visualize** geofence boundaries on the map

### Vehicle Management

- **Register** new vehicles with driver information
- **View** all registered vehicles
- **Update** vehicle locations manually or via API
- **Track** vehicle positions on the map
- **Monitor** vehicle status (active/inactive)

### Alert Management

- **Configure** alert rules for specific geofences
- **Set** event types (entry, exit, or both)
- **Apply** rules to specific vehicles or all vehicles
- **Monitor** real-time alerts as they occur
- **View** detailed alert feed with timestamps

### Violation History

- **Filter** violations by vehicle, geofence, and date range
- **View** detailed information about each event
- **Track** entry/exit patterns
- **Export** violation data (feature can be added)

### Real-time Alerts

- **Receive** instant notifications when vehicles enter/exit geofences
- **Color-coded** alerts (red for restricted zones, green for delivery zones)
- **View** alert details including location and timestamp
- **Audio/visual** notifications for critical events

---

## Architecture Overview

### Layered Architecture

#### 1. **Domain Layer** (`internal/domain/`)

- Business logic interfaces
- Domain models (Geofence, Vehicle, Alert)
- Error definitions
- No external dependencies

#### 2. **Repository Layer** (`internal/repositories/`)

- Data persistence abstraction
- Database queries implementation
- Transaction management
- Dependency: Database, Domain

#### 3. **Service Layer** (`internal/services/`)

- Business logic implementation
- Orchestration of repositories
- Data transformation
- Domain rule enforcement

#### 4. **Handler Layer** (`internal/handlers/`)

- HTTP request handling
- Input validation
- Response formatting
- Dependency: Services

### SOLID Principles Applied

#### Single Responsibility Principle (SRP)

- Each service has one reason to change
- Separate concerns (handlers, services, repositories)

#### Open/Closed Principle (OCP)

- Open for extension via interfaces
- Closed for modification (interfaces define contracts)

#### Liskov Substitution Principle (LSP)

- Repositories implement consistent interfaces
- Services can be swapped for testing

#### Interface Segregation Principle (ISP)

- Domain interfaces are focused and minimal
- Clients depend only on needed methods

#### Dependency Inversion Principle (DIP)

- Services depend on interfaces, not concrete types
- Repositories injected into services
- Handlers receive initialized services

### Component Diagram

```
┌─────────────────────────────────────┐
│         Frontend (React)             │
│  - Dashboard, Forms, Maps, Alerts    │
└──────────────┬──────────────────────┘
               │
        HTTP / WebSocket
               │
┌──────────────▼──────────────────────┐
│        API Handlers Layer            │
│  - GeofenceHandler                  │
│  - VehicleHandler                   │
│  - AlertHandler                     │
│  - WebSocketHandler                 │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│        Service Layer                 │
│  - GeofenceService                  │
│  - VehicleService                   │
│  - AlertService                     │
│  - Point-in-polygon detection       │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      Repository Layer                │
│  - GeofenceRepository               │
│  - VehicleRepository                │
│  - VehicleLocationRepository        │
│  - AlertRepository                  │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      PostgreSQL Database             │
│  - Geofences, Vehicles, Events      │
│  - Locations, Alert Rules           │
└──────────────────────────────────────┘
```

### Data Flow: Location Update

```
POST /vehicles/location
         │
         ▼
  VehicleHandler
         │
         ├─▶ UpdateVehicleLocation (VehicleService)
         │         │
         │         ├─▶ SaveLocation (VehicleLocationRepository)
         │         │
         │         └─▶ GetCurrentGeofences
         │                  │
         │                  └─▶ Point-in-Polygon Check
         │
         ├─▶ CheckAndTriggerAlerts (AlertService)
         │         │
         │         ├─▶ Get matching alert rules
         │         │
         │         ├─▶ SaveEvent (AlertRepository)
         │         │
         │         └─▶ PublishAlert (via WebSocket)
         │                  │
         │                  └─▶ All connected clients receive alert
         │
         ▼
  HTTP Response
```

---

##Testing

### Unit Tests

```bash
cd backend

# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestGeofenceService ./internal/services/...
```

### API Testing with Postman

1. **Import Collection**
   - Open Postman
   - Click "Import"
   - Import `postman_collection.json` (create or download)

2. **Create Test Environment**
   - Set `base_url` = `http://localhost:8080`
   - Set `vehicle_id`, `geofence_id` from responses

3. **Run Test Sequences**
   - Create geofence
   - Register vehicle
   - Update location
   - Configure alert
   - Verify violation history

### Manual Testing

```bash
# Test geofence creation
curl -X POST http://localhost:8080/geofences \
  -H "Content-Type: application/json" \
  -d @geofence-test.json

# Test vehicle registration
curl -X POST http://localhost:8080/vehicles \
  -H "Content-Type: application/json" \
  -d @vehicle-test.json

# Test location update
curl -X POST http://localhost:8080/vehicles/location \
  -H "Content-Type: application/json" \
  -d @location-test.json
```

---

## Troubleshooting

### Database Connection Issues

**Error**: `connection refused`

**Solutions**:

```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Check database credentials in .env
cat .env | grep DB_

# Test connection
psql -U postgres -h localhost -d geofence -c "SELECT 1;"
```

### Backend Won't Start

**Error**: `port 8080 already in use`

**Solutions**:

```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>

# Or use different port
SERVER_PORT=8081 go run ./cmd/main.go
```

### Frontend API Connection Issues

**Error**: `Failed to fetch from localhost:8080`

**Solutions**:

```bash
# Check REACT_APP_API_URL in .env.local
cat frontend/.env.local

# Check CORS is enabled in backend
# Should see Access-Control-Allow headers in response

# Check backend is running
curl http://localhost:8080/geofences
```

### WebSocket Connection Issues

**Error**: `WebSocket connection failed`

**Solutions**:

```bash
# Check WebSocket endpoint is accessible
curl --include \
  --no-buffer \
  --header "Connection: Upgrade" \
  --header "Upgrade: websocket" \
  --header "Sec-WebSocket-Version: 13" \
  --header "Sec-WebSocket-Key: SGVsbG8=" \
  http://localhost:8080/ws/alerts

# Check firewall allows WebSocket
# Check reverse proxy (if any) supports WebSocket upgrades
```

### Performance Issues

**Slow Geofence Detection**:

```bash
# Check geofence complexity (number of points)
# Simplify polygon if needed
# Add database indexes (auto-created)

# Monitor backend logs
docker-compose logs backend | grep -i error
```

### Docker Issues

**Containers not starting**:

```bash
# Check logs
docker-compose logs

# Restart services
docker-compose restart

# Full rebuild
docker-compose down
docker-compose up --build
```

---

## Production Checklist

- [ ] Set strong database passwords
- [ ] Enable HTTPS/TLS
- [ ] Implement authentication (JWT/API keys)
- [ ] Configure rate limiting
- [ ] Enable CORS properly
- [ ] Setup monitoring and alerting
- [ ] Configure backups
- [ ] Enable database replication
- [ ] Use environment variables for secrets
- [ ] Test load handling
- [ ] Setup log aggregation
- [ ] Document deployment process
- [ ] Create runbooks for operations
- [ ] Plan disaster recovery

---

## Support & Resources

- **Go Documentation**: https://golang.org/doc/
- **React Documentation**: https://react.dev/
- **PostgreSQL Documentation**: https://www.postgresql.org/docs/
- **Leaflet Documentation**: https://leafletjs.com/
- **Docker Documentation**: https://docs.docker.com/

---

## License

This project is provided as-is for assessment purposes.
