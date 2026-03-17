import { useCallback } from "react";
import { api } from "../services/api";

export const useGeofenceApi = () => {
  const createGeofence = useCallback(async (data) => {
    return await api.createGeofence(data);
  }, []);

  const getGeofences = useCallback(async (category) => {
    return await api.getGeofences(category);
  }, []);

  return { createGeofence, getGeofences };
};

export const useVehicleApi = () => {
  const registerVehicle = useCallback(async (data) => {
    return await api.registerVehicle(data);
  }, []);

  const getVehicles = useCallback(async () => {
    return await api.getVehicles();
  }, []);

  const updateLocation = useCallback(async (data) => {
    return await api.updateVehicleLocation(data);
  }, []);

  const getLocation = useCallback(async (vehicleId) => {
    return await api.getVehicleLocation(vehicleId);
  }, []);

  return { registerVehicle, getVehicles, updateLocation, getLocation };
};

export const useAlertApi = () => {
  const configureAlert = useCallback(async (data) => {
    return await api.configureAlert(data);
  }, []);

  const getAlerts = useCallback(async (filters) => {
    return await api.getAlerts(filters);
  }, []);

  const getViolationHistory = useCallback(async (filters) => {
    return await api.getViolationHistory(filters);
  }, []);

  return { configureAlert, getAlerts, getViolationHistory };
};
