import React, { useEffect, useState } from "react";
import {
  useGeofenceStore,
  useVehicleStore,
  useWebSocketStore,
} from "../services/store";
import { useGeofenceApi, useVehicleApi } from "../hooks/useApi";
import MapComponent from "../components/MapComponent";
import AlertsFeed from "../components/AlertsFeed";

const DashboardPage = () => {
  const { geofences = [], setGeofences } = useGeofenceStore();
  const { vehicles = [], setVehicles, locations = {} } = useVehicleStore();
  const { connected } = useWebSocketStore();
  const { getGeofences } = useGeofenceApi();
  const { getVehicles } = useVehicleApi();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      const [geoRes, vehRes] = await Promise.all([
        getGeofences(),
        getVehicles(),
      ]);

      if (geoRes.geofences) {
        setGeofences(geoRes.geofences);
      }

      if (vehRes.vehicles) {
        setVehicles(vehRes.vehicles);
      }
    } catch (error) {
      console.error("Failed to load data", error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="flex justify-between items-center mb-4">
          <h1 className="text-3xl font-bold">Dashboard</h1>
          <div className="flex items-center gap-4">
            <div
              className={`flex items-center gap-2 px-4 py-2 rounded-lg ${connected ? "bg-success text-white" : "bg-danger text-white"}`}
            >
              <span
                className={`inline-block w-3 h-3 rounded-full ${connected ? "bg-white" : "opacity-50"}`}
              ></span>
              {connected ? "Connected" : "Disconnected"}
            </div>
            <button
              onClick={loadData}
              disabled={loading}
              className="bg-primary text-white px-4 py-2 rounded hover:bg-opacity-90 transition disabled:opacity-50"
            >
              {loading ? "Refreshing..." : "Refresh"}
            </button>
          </div>
        </div>

        {/* Statistics */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="bg-gradient-to-br from-blue-500 to-blue-600 text-white rounded-lg p-4">
            <h3 className="text-sm font-medium opacity-80">Geofences</h3>
            <p className="text-3xl font-bold">{geofences.length}</p>
          </div>
          <div className="bg-gradient-to-br from-green-500 to-green-600 text-white rounded-lg p-4">
            <h3 className="text-sm font-medium opacity-80">Vehicles</h3>
            <p className="text-3xl font-bold">{vehicles.length}</p>
          </div>
          <div className="bg-gradient-to-br from-purple-500 to-purple-600 text-white rounded-lg p-4">
            <h3 className="text-sm font-medium opacity-80">Active Tracking</h3>
            <p className="text-3xl font-bold">
              {Object.keys(locations).length}
            </p>
          </div>
        </div>
      </div>

      {/* Map */}
      <div
        className="bg-white rounded-lg shadow-md overflow-hidden"
        style={{ height: "500px" }}
      >
        <MapComponent geofences={geofences} vehicles={locations} />
      </div>

      {/* Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          {/* Geofences List */}
          <div className="bg-white rounded-lg shadow-md p-6 mb-6">
            <h2 className="text-xl font-bold mb-4">Recent Geofences</h2>
            {geofences.length === 0 ? (
              <p className="text-gray-500">No geofences created yet</p>
            ) : (
              <div className="space-y-3 max-h-64 overflow-y-auto">
                {geofences.slice(0, 5).map((g) => (
                  <div
                    key={g.id}
                    className="border rounded-lg p-3 hover:bg-gray-50"
                  >
                    <h3 className="font-semibold">{g.name}</h3>
                    <p className="text-sm text-gray-600">{g.description}</p>
                    <div className="flex justify-between items-center mt-2">
                      <span
                        className={`text-xs px-2 py-1 rounded ${
                          g.category === "restricted_zone"
                            ? "bg-danger text-white"
                            : "bg-secondary text-white"
                        }`}
                      >
                        {g.category}
                      </span>
                      <small className="text-gray-500">
                        {new Date(g.created_at).toLocaleDateString()}
                      </small>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Vehicles List */}
          <div className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-xl font-bold mb-4">Registered Vehicles</h2>
            {vehicles.length === 0 ? (
              <p className="text-gray-500">No vehicles registered yet</p>
            ) : (
              <div className="space-y-3 max-h-64 overflow-y-auto">
                {vehicles.slice(0, 5).map((v) => (
                  <div
                    key={v.id}
                    className="border rounded-lg p-3 hover:bg-gray-50"
                  >
                    <h3 className="font-semibold">{v.vehicle_number}</h3>
                    <p className="text-sm text-gray-600">
                      Driver: {v.driver_name}
                    </p>
                    <div className="flex justify-between items-center mt-2">
                      <span className="text-xs px-2 py-1 rounded bg-success text-white">
                        {v.vehicle_type}
                      </span>
                      <small className="text-gray-500">{v.phone}</small>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Alerts Feed */}
        <div className="lg:col-span-1">
          <AlertsFeed />
        </div>
      </div>
    </div>
  );
};

export default DashboardPage;
