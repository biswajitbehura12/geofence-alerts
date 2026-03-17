import { useCallback, useEffect, useState } from "react";
import { useGeofenceStore } from "../services/store";
import { useGeofenceApi } from "../hooks/useApi";
import GeofenceForm from "../components/GeofenceForm";
import MapComponent from "../components/MapComponent";

const GeofencesPage = () => {
  const { geofences = [], setGeofences } = useGeofenceStore();
  const { getGeofences } = useGeofenceApi();
  const [loading, setLoading] = useState(false);
  const [selectedLocation, setSelectedLocation] = useState(null);

  const loadGeofences = useCallback(async () => {
    setLoading(true);
    try {
      const result = await getGeofences();
      if (result.geofences) {
        setGeofences(result.geofences);
      }
    } catch (error) {
      console.error("Failed to load geofences", error);
    } finally {
      setLoading(false);
    }
  }, [getGeofences, setGeofences]);

  useEffect(() => {
    loadGeofences();
  }, [loadGeofences]);

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold">Geofence Management</h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Form */}
        <div className="lg:col-span-1">
          <GeofenceForm onSuccess={loadGeofences} />
        </div>

        {/* Map and List */}
        <div className="lg:col-span-2 space-y-6">
          {/* Map */}
          <div
            className="bg-white rounded-lg shadow-md overflow-hidden"
            style={{ height: "400px" }}
          >
            <MapComponent
              geofences={geofences}
              onMapClick={setSelectedLocation}
              selectedLocation={selectedLocation}
            />
          </div>

          {/* Geofences List */}
          <div className="bg-white rounded-lg shadow-md p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-2xl font-bold">All Geofences</h2>
              <button
                onClick={loadGeofences}
                disabled={loading}
                className="bg-primary text-white px-4 py-2 rounded hover:bg-opacity-90 transition disabled:opacity-50"
              >
                {loading ? "Loading..." : "Refresh"}
              </button>
            </div>

            {geofences.length === 0 ? (
              <p className="text-gray-500 text-center py-8">
                No geofences created yet
              </p>
            ) : (
              <div className="grid grid-cols-1 gap-4 max-h-96 overflow-y-auto">
                {geofences.map((geofence) => (
                  <div
                    key={geofence.id}
                    className="border rounded-lg p-4 hover:shadow-md transition cursor-pointer"
                    onClick={() =>
                      setSelectedLocation({
                        lat: geofence.coordinates[0][0],
                        lng: geofence.coordinates[0][1],
                      })
                    }
                  >
                    <div className="flex justify-between items-start mb-2">
                      <h3 className="font-bold text-lg">{geofence.name}</h3>
                      <span
                        className={`text-xs px-3 py-1 rounded text-white ${
                          geofence.category === "restricted_zone"
                            ? "bg-danger"
                            : geofence.category === "delivery_zone"
                              ? "bg-success"
                              : "bg-secondary"
                        }`}
                      >
                        {geofence.category}
                      </span>
                    </div>
                    <p className="text-sm text-gray-600 mb-2">
                      {geofence.description}
                    </p>
                    <div className="text-xs text-gray-500 grid grid-cols-2">
                      <div>Points: {geofence.coordinates.length}</div>
                      <div>
                        {new Date(geofence.created_at).toLocaleDateString()}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default GeofencesPage;
