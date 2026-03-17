# Geofence & Vehicle Tracking System

A comprehensive, production-grade full-stack application for geofencing and real-time vehicle tracking with industrial-grade code structure following SOLID principles.

## 🎯 Features

### Core Functionality

- ✅ **Geofence Management** - Create and manage virtual boundaries
- ✅ **Vehicle Registration** - Register and track vehicles
- ✅ **Real-time Location Tracking** - Update vehicle positions
- ✅ **Geofence Detection** - Accurate point-in-polygon detection
- ✅ **Alert Configuration** - Define custom alert rules
- ✅ **Real-time Alerts** - WebSocket-based instant notifications
- ✅ **Violation History** - Track all geofence events
- ✅ **Interactive Dashboard** - Beautiful, responsive UI with maps

### Technical Highlights

- 🏗️ **Layered Architecture** - Domain, Repository, Service, Handler layers
- 🔐 **SOLID Principles** - Single Responsibility, Open/Closed, Liskov Substitution, Interface Segregation, Dependency Inversion
- 📊 **PostgreSQL** - Robust relational database with proper indexing
- 🚀 **Go Backend** - Fast, concurrent processing with goroutines
- ⚛️ **React Frontend** - Modern, reusable components with hooks
- 🗺️ **Leaflet Maps** - Interactive map visualization
- 🔔 **WebSocket** - Real-time bidirectional communication
- 🐳 **Docker** - Complete containerization with docker-compose
- 📝 **Comprehensive Documentation** - Setup guides, API docs, troubleshooting

## 🛠️ Technology Stack

### Backend

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP)
- **Database**: PostgreSQL 13+
- **Real-time**: Gorilla WebSocket
- **Architecture**: Layered with DI

### Frontend

- **Framework**: React 18+
- **State Management**: Zustand
- **Styling**: Tailwind CSS
- **Maps**: Leaflet.js
- **Notifications**: React Toastify
- **Http Client**: Fetch API

### Infrastructure

- **Containerization**: Docker & Docker Compose
- **Database**: PostgreSQL
- **Deployment**: Cloud-ready (AWS, GCP, Azure)

## 📁 Project Structure

```
geofence/
├── backend/
│   ├── cmd/
│   │   └── main.go                 # Entry point
│   ├── config/
│   │   └── config.go               # Configuration management
│   ├── internal/
│   │   ├── domain/                 # Domain models & interfaces
│   │   │   ├── geofence.go
│   │   │   ├── vehicle.go
│   │   │   ├── alert.go
│   │   │   └── errors.go
│   │   ├── repositories/           # Data persistence
│   │   │   ├── geofence_repository.go
│   │   │   ├── vehicle_repository.go
│   │   │   └── alert_repository.go
│   │   ├── services/               # Business logic
│   │   │   ├── geofence_service.go
│   │   │   ├── vehicle_service.go
│   │   │   ├── alert_service.go
│   │   │   └── time_helper.go
│   │   ├── handlers/               # HTTP handlers
│   │   │   ├── geofence_handler.go
│   │   │   ├── vehicle_handler.go
│   │   │   ├── alert_handler.go
│   │   │   └── websocket_handler.go
│   │   └── models/                 # Request/Response models
│   │       ├── models.go
│   │       └── requests_responses.go
│   ├── migrations/                 # Database migrations
│   ├── docker/
│   │   └── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── .env.example
│
├── frontend/
│   ├── src/
│   │   ├── components/             # Reusable React components
│   │   │   ├── Navbar.jsx
│   │   │   ├── MapComponent.jsx
│   │   │   ├── GeofenceForm.jsx
│   │   │   ├── VehicleForm.jsx
│   │   │   ├── AlertConfiguration.jsx
│   │   │   ├── AlertsFeed.jsx
│   │   │   ├── LocationUpdater.jsx
│   │   │   └── ViolationHistory.jsx
│   │   ├── pages/                  # Page components
│   │   │   ├── Dashboard.jsx
│   │   │   ├── Geofences.jsx
│   │   │   ├── Vehicles.jsx
│   │   │   └── Alerts.jsx
│   │   ├── services/               # API & Business logic
│   │   │   ├── api.js
│   │   │   └── store.js
│   │   ├── hooks/                  # Custom React hooks
│   │   │   ├── useApi.js
│   │   │   └── useWebSocket.js
│   │   ├── styles/                 # Global styles
│   │   │   └── globals.css
│   │   ├── App.jsx
│   │   └── index.jsx
│   ├── public/
│   │   └── index.html
│   ├── Dockerfile
│   ├── package.json
│   ├── tailwind.config.js
│   ├── postcss.config.js
│   └── .env.example
│
├── docker-compose.yml              # Multi-container setup
├── SETUP.md                        # Comprehensive setup guide
├── README.md                       # This file
└── .env.example
```

## 🚀 Quick Start

### Using Docker Compose (Recommended)

```bash
# Clone repository
git clone <repository-url>
cd geofence

# Copy environment file
cp .env.example .env

# Start all services
docker-compose up --build

# Application will be available at:
# Frontend: http://localhost:3000
# Backend: http://localhost:8080
# Database: localhost:5432
```

### Local Development

See [SETUP.md](./SETUP.md) for detailed local setup instructions.

```bash
# Setup backend
cd backend
go mod download
export DB_HOST=localhost
go run ./cmd/main.go

# Setup frontend (in new terminal)
cd frontend
npm install
npm start
```

## 📡 API Endpoints

### Geofence Management

- `POST /geofences` - Create geofence
- `GET /geofences` - Get all geofences

### Vehicle Management

- `POST /vehicles` - Register vehicle
- `GET /vehicles` - Get all vehicles
- `POST /vehicles/location` - Update vehicle location
- `GET /vehicles/location/{id}` - Get vehicle location

### Alert Management

- `POST /alerts/configure` - Configure alert rule
- `GET /alerts` - Get all alert rules
- `GET /violations/history` - Get violation history

### Real-time

- `WebSocket /ws/alerts` - Real-time alert stream

See [SETUP.md](./SETUP.md) for detailed API documentation.

## 🏛️ Architecture Highlights

### Layered Architecture

```
┌─────────────────────────────────┐
│       HTTP Handlers             │
│  - Request parsing & response   │
└────────────┬────────────────────┘
             │
┌────────────▼────────────────────┐
│       Service Layer             │
│  - Business logic               │
│  - Orchestration                │
│  - Validation                   │
└────────────┬────────────────────┘
             │
┌────────────▼────────────────────┐
│      Repository Layer           │
│  - Data persistence             │
│  - Database queries             │
└────────────┬────────────────────┘
             │
┌────────────▼────────────────────┐
│      Domain Layer               │
│  - Models & Interfaces          │
│  - Business rules               │
└────────────────────────────────┘
```

### SOLID Principles

✅ **Single Responsibility** - Each layer has one reason to change
✅ **Open/Closed** - Open for extension via interfaces
✅ **Liskov Substitution** - Repositories implement consistent contracts
✅ **Interface Segregation** - Focused, minimal interfaces
✅ **Dependency Inversion** - Depend on abstractions, not concrete types

## 🔍 Key Algorithms

### Point-in-Polygon Detection

Uses ray casting algorithm for accurate geofence membership checking:

- Time complexity: O(n) where n is number of polygon vertices
- Works with concave and convex polygons
- Handles edge cases (point on boundary, multiple intersections)

### Real-time Alert Processing

Asynchronous alert processing with WebSocket delivery:

- Non-blocking event processing
- In-memory subscriber management
- Automatic connection recovery

## 📊 Database Schema

### Geofences Table

```sql
- id (VARCHAR): Primary key
- name (VARCHAR): Geofence name
- description (TEXT): Description
- coordinates (FLOAT8[]): Polygon boundary points
- category (VARCHAR): delivery_zone, restricted_zone, toll_zone, customer_area
- created_at, updated_at: Timestamps
```

### Vehicles Table

```sql
- id (VARCHAR): Primary key
- vehicle_number (VARCHAR): Unique registration
- driver_name (VARCHAR): Driver name
- vehicle_type (VARCHAR): Type of vehicle
- phone (VARCHAR): Contact number
- created_at, updated_at: Timestamps
```

### Additional Tables

- `vehicle_locations` - Location history with timestamps
- `alert_rules` - Alert configuration rules
- `geofence_events` - Violation history

## 🧪 Testing

### Backend Tests

```bash
cd backend
go test ./...
go test -cover ./...
```

### Frontend Tests

```bash
cd frontend
npm test
```

### API Testing

Use Postman collection or curl commands (see SETUP.md)

## 🐳 Docker Deployment

### Build Images

```bash
# Backend
docker build -f backend/docker/Dockerfile -t geofence-backend:latest ./backend

# Frontend
docker build -t geofence-frontend:latest ./frontend
```

### Run with Docker Compose

```bash
docker-compose up -d
```

## 📈 Performance Optimization

### Backend

- ✅ Connection pooling (PostgreSQL)
- ✅ Query optimization with indexes
- ✅ Concurrent request handling
- ✅ Efficient point-in-polygon algorithm
- ✅ Asynchronous alert processing

### Frontend

- ✅ Component memoization
- ✅ State management with Zustand
- ✅ Lazy loading of components
- ✅ Optimized re-renders
- ✅ Efficient WebSocket handling

## 🔒 Security Considerations

For production deployment:

- [ ] Implement JWT authentication
- [ ] Add rate limiting
- [ ] Enable HTTPS/TLS
- [ ] Use environment variables for secrets
- [ ] Implement CORS properly
- [ ] Add input validation
- [ ] Use parameterized queries
- [ ] Add request logging
- [ ] Implement audit trails

## 🚀 Deployment

### AWS EC2

```bash
# Backend deployment
# Frontend on Vercel/Netlify
# Database on RDS PostgreSQL
```

### Google Cloud

```bash
# Cloud Run for backend
# Cloud Storage for frontend
# Cloud SQL for database
```

### Azure

```bash
# App Service for backend
# Static Web Apps for frontend
# Azure Database for PostgreSQL
```

See SETUP.md for detailed deployment instructions.

## 📚 Documentation

- [SETUP.md](./SETUP.md) - Complete setup and configuration guide
- [API Documentation](#-api-endpoints) - API endpoint reference
- [Architecture Guide](#-architecture-highlights) - System design
- [Troubleshooting](./SETUP.md#troubleshooting) - Common issues and solutions

## 🤝 Contributing

This is an assessment project. Please follow:

- Go code conventions (gofmt, golint)
- React best practices (hooks, functional components)
- SOLID principles in design
- Comprehensive error handling
- Clear naming conventions

## 📝 License

Assessment project - Provided as-is for evaluation.

## 📞 Support

For issues or questions about the implementation:

1. Check SETUP.md Troubleshooting section
2. Review API documentation
3. Check logs: `docker-compose logs -f`
4. Verify environment configuration

---

**Built with ❤️ following industry best practices and SOLID principles**
