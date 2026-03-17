import { useState } from "react";
import { useVehicleStore, useGeofenceStore } from "../services/store";
import { useVehicleApi, useAlertApi } from "../hooks/useApi";
import { toast } from "react-toastify";

const LocationUpdater = () => {
  const { vehicles = [] } = useVehicleStore();
  const { updateLocation } = useVehicleApi();
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    vehicle_id: "",
    latitude: "",
    longitude: "",
    timestamp: new Date().toISOString(),
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (
      !formData.vehicle_id ||
      formData.latitude === "" ||
      formData.longitude === ""
    ) {
      toast.error("All fields are required");
      return;
    }

    const lat = parseFloat(formData.latitude);
    const lng = parseFloat(formData.longitude);

    if (lat < -90 || lat > 90 || lng < -180 || lng > 180) {
      toast.error("Invalid coordinates");
      return;
    }

    setLoading(true);
    try {
      const result = await updateLocation({
        vehicle_id: formData.vehicle_id,
        latitude: lat,
        longitude: lng,
        timestamp: formData.timestamp,
      });

      if (result.error) {
        toast.error(result.error);
      } else {
        toast.success("Location updated successfully");
        setFormData({
          vehicle_id: "",
          latitude: "",
          longitude: "",
          timestamp: new Date().toISOString(),
        });
      }
    } catch (error) {
      toast.error("Failed to update location");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6 mb-6">
      <h2 className="text-2xl font-bold mb-4">Update Vehicle Location</h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-1">Vehicle *</label>
          <select
            name="vehicle_id"
            value={formData.vehicle_id}
            onChange={handleChange}
            className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
          >
            <option value="">Select a vehicle</option>
            {vehicles.map((v) => (
              <option key={v.id} value={v.id}>
                {v.vehicle_number} - {v.driver_name}
              </option>
            ))}
          </select>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-1">Latitude *</label>
            <input
              type="number"
              step="0.0001"
              min="-90"
              max="90"
              name="latitude"
              value={formData.latitude}
              onChange={handleChange}
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="-90 to 90"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">
              Longitude *
            </label>
            <input
              type="number"
              step="0.0001"
              min="-180"
              max="180"
              name="longitude"
              value={formData.longitude}
              onChange={handleChange}
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="-180 to 180"
            />
          </div>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full bg-primary text-white py-2 rounded hover:bg-opacity-90 transition disabled:opacity-50"
        >
          {loading ? "Updating..." : "Update Location"}
        </button>
      </form>
    </div>
  );
};

export default LocationUpdater;
