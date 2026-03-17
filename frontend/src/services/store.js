import create from "zustand";

// Geofence store
export const useGeofenceStore = create((set) => ({
  geofences: [],
  selectedGeofence: null,
  loading: false,
  error: null,

  setGeofences: (geofences) => set({ geofences }),
  setSelectedGeofence: (geofence) => set({ selectedGeofence: geofence }),
  setLoading: (loading) => set({ loading }),
  setError: (error) => set({ error }),
  addGeofence: (geofence) =>
    set((state) => ({
      geofences: [...state.geofences, geofence],
    })),
}));

// Vehicle store
export const useVehicleStore = create((set) => ({
  vehicles: [],
  locations: {},
  loading: false,
  error: null,

  setVehicles: (vehicles) => set({ vehicles }),
  setLocations: (locations) => set({ locations }),
  updateVehicleLocation: (vehicleId, location) =>
    set((state) => ({
      locations: { ...state.locations, [vehicleId]: location },
    })),
  setLoading: (loading) => set({ loading }),
  setError: (error) => set({ error }),
  addVehicle: (vehicle) =>
    set((state) => ({
      vehicles: [...state.vehicles, vehicle],
    })),
}));

// Alert store
export const useAlertStore = create((set) => ({
  alerts: [],
  recentAlerts: [],
  violations: [],
  loading: false,

  setAlerts: (alerts) => set({ alerts }),
  setViolations: (violations) => set({ violations }),
  addAlert: (alert) =>
    set((state) => ({
      recentAlerts: [alert, ...state.recentAlerts.slice(0, 99)],
    })),
  setLoading: (loading) => set({ loading }),
}));

// WebSocket store
export const useWebSocketStore = create((set) => ({
  connected: false,
  setConnected: (connected) => set({ connected }),
}));
