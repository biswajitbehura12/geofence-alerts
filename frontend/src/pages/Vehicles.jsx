import React, { useEffect, useState } from "react";
import { useVehicleStore, useGeofenceStore } from "../services/store";
import { useVehicleApi } from "../hooks/useApi";
import VehicleForm from "../components/VehicleForm";
import LocationUpdater from "../components/LocationUpdater";
import MapComponent from "../components/MapComponent";

const VehiclesPage = () => {
  const { vehicles = [], setVehicles, locations = {} } = useVehicleStore();
  const { geofences = [] } = useGeofenceStore();
  const { getVehicles } = useVehicleApi();
  const [loading, setLoading] = useState(false);
  const [selectedLocation, setSelectedLocation] = useState(null);

  useEffect(() => {
    loadVehicles();
  }, []);

  const loadVehicles = async () => {
    setLoading(true);
    try {
      const result = await getVehicles();
      if (result.vehicles) {
        setVehicles(result.vehicles);
      }
    } catch (error) {
      console.error("Failed to load vehicles", error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold">Vehicle Management</h1>

      {/* Forms */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <VehicleForm onSuccess={loadVehicles} />
        <LocationUpdater />
      </div>

      {/* Map and List */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Map */}
        <div
          className="lg:col-span-2 bg-white rounded-lg shadow-md overflow-hidden"
          style={{ height: "400px" }}
        >
          <MapComponent
            geofences={geofences}
            vehicles={locations}
            onMapClick={setSelectedLocation}
            selectedLocation={selectedLocation}
          />
        </div>

        {/* Quick Stats */}
        <div className="space-y-4">
          <div className="bg-white rounded-lg shadow-md p-6">
            <h3 className="font-bold text-lg mb-2">Total Vehicles</h3>
            <p className="text-4xl font-bold text-primary">{vehicles.length}</p>
          </div>

          <div className="bg-white rounded-lg shadow-md p-6">
            <h3 className="font-bold text-lg mb-2">Tracked Vehicles</h3>
            <p className="text-4xl font-bold text-success">
              {Object.keys(locations).length}
            </p>
          </div>
        </div>
      </div>

      {/* Vehicles List */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-2xl font-bold">All Vehicles</h2>
          <button
            onClick={loadVehicles}
            disabled={loading}
            className="bg-primary text-white px-4 py-2 rounded hover:bg-opacity-90 transition disabled:opacity-50"
          >
            {loading ? "Loading..." : "Refresh"}
          </button>
        </div>

        {vehicles.length === 0 ? (
          <p className="text-gray-500 text-center py-8">
            No vehicles registered yet
          </p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b">
                  <th className="text-left py-3 px-4">Vehicle Number</th>
                  <th className="text-left py-3 px-4">Driver</th>
                  <th className="text-left py-3 px-4">Type</th>
                  <th className="text-left py-3 px-4">Phone</th>
                  <th className="text-left py-3 px-4">Status</th>
                  <th className="text-left py-3 px-4">Tracking</th>
                  <th className="text-left py-3 px-4">Created</th>
                </tr>
              </thead>
              <tbody>
                {vehicles.map((vehicle) => {
                  const tracking = locations[vehicle.id];
                  return (
                    <tr key={vehicle.id} className="border-b hover:bg-gray-50">
                      <td className="py-3 px-4 font-semibold">
                        {vehicle.vehicle_number}
                      </td>
                      <td className="py-3 px-4">{vehicle.driver_name}</td>
                      <td className="py-3 px-4">{vehicle.vehicle_type}</td>
                      <td className="py-3 px-4 text-xs">{vehicle.phone}</td>
                      <td className="py-3 px-4">
                        <span className="bg-success text-white px-2 py-1 rounded text-xs">
                          {vehicle.status}
                        </span>
                      </td>
                      <td className="py-3 px-4">
                        {tracking ? (
                          <span className="bg-warning text-black px-2 py-1 rounded text-xs">
                            🔴 Active
                          </span>
                        ) : (
                          <span className="bg-gray-300 text-gray-700 px-2 py-1 rounded text-xs">
                            Inactive
                          </span>
                        )}
                      </td>
                      <td className="py-3 px-4 text-xs">
                        {new Date(vehicle.created_at).toLocaleDateString()}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
};

export default VehiclesPage;
