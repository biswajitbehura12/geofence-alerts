// API service for backend communication
const API_BASE_URL = process.env.REACT_APP_API_URL || "http://localhost:8080";

export const api = {
  // Geofences
  createGeofence: async (data) => {
    const response = await fetch(`${API_BASE_URL}/geofences`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    return response.json();
  },

  getGeofences: async (category) => {
    let url = `${API_BASE_URL}/geofences`;
    if (category) url += `?category=${category}`;
    const response = await fetch(url);
    return response.json();
  },

  // Vehicles
  registerVehicle: async (data) => {
    const response = await fetch(`${API_BASE_URL}/vehicles`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    return response.json();
  },

  getVehicles: async () => {
    const response = await fetch(`${API_BASE_URL}/vehicles`);
    return response.json();
  },

  updateVehicleLocation: async (data) => {
    const response = await fetch(`${API_BASE_URL}/vehicles/location`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    return response.json();
  },

  getVehicleLocation: async (vehicleId) => {
    const response = await fetch(
      `${API_BASE_URL}/vehicles/location/${vehicleId}`,
    );
    return response.json();
  },

  // Alerts
  configureAlert: async (data) => {
    const response = await fetch(`${API_BASE_URL}/alerts/configure`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    return response.json();
  },

  getAlerts: async (filters = {}) => {
    let url = `${API_BASE_URL}/alerts`;
    const params = new URLSearchParams();
    if (filters.geofenceId) params.append("geofence_id", filters.geofenceId);
    if (filters.vehicleId) params.append("vehicle_id", filters.vehicleId);
    if (params.toString()) url += `?${params.toString()}`;
    const response = await fetch(url);
    return response.json();
  },

  getViolationHistory: async (filters = {}) => {
    let url = `${API_BASE_URL}/violations/history`;
    const params = new URLSearchParams();
    if (filters.vehicleId) params.append("vehicle_id", filters.vehicleId);
    if (filters.geofenceId) params.append("geofence_id", filters.geofenceId);
    if (filters.startDate) params.append("start_date", filters.startDate);
    if (filters.endDate) params.append("end_date", filters.endDate);
    if (filters.limit) params.append("limit", filters.limit);
    if (filters.offset) params.append("offset", filters.offset);
    if (params.toString()) url += `?${params.toString()}`;
    const response = await fetch(url);
    return response.json();
  },
};

// WebSocket service for real-time alerts
export const createWebSocketConnection = () => {
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  const wsURL =
    process.env.REACT_APP_WS_URL ||
    `${protocol}//${window.location.host.split(":")[0]}:8080/ws/alerts`;
  return new WebSocket(wsURL);
};
