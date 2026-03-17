import React, { useState } from "react";
import { useGeofenceStore } from "../services/store";
import { useGeofenceApi } from "../hooks/useApi";
import { toast } from "react-toastify";

const GeofenceForm = ({ onSuccess }) => {
  const { addGeofence } = useGeofenceStore();
  const { createGeofence } = useGeofenceApi();
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    name: "",
    description: "",
    category: "delivery_zone",
    coordinates: [],
  });

  const handleAddCoordinate = (e) => {
    e.preventDefault();
    const lat = parseFloat(document.getElementById("lat").value);
    const lng = parseFloat(document.getElementById("lng").value);

    if (isNaN(lat) || isNaN(lng)) {
      toast.error("Invalid coordinates");
      return;
    }

    if (lat < -90 || lat > 90 || lng < -180 || lng > 180) {
      toast.error("Coordinates out of valid range");
      return;
    }

    setFormData({
      ...formData,
      coordinates: [...formData.coordinates, [lat, lng]],
    });

    document.getElementById("lat").value = "";
    document.getElementById("lng").value = "";
  };

  const handleRemoveCoordinate = (index) => {
    setFormData({
      ...formData,
      coordinates: formData.coordinates.filter((_, i) => i !== index),
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!formData.name) {
      toast.error("Name is required");
      return;
    }

    if (formData.coordinates.length < 3) {
      toast.error("At least 3 coordinates required");
      return;
    }

    // Close polygon
    const closedCoords = [...formData.coordinates];
    if (
      closedCoords[0][0] !== closedCoords[closedCoords.length - 1][0] ||
      closedCoords[0][1] !== closedCoords[closedCoords.length - 1][1]
    ) {
      closedCoords.push(closedCoords[0]);
    }

    setLoading(true);
    try {
      const result = await createGeofence({
        ...formData,
        coordinates: closedCoords,
      });

      if (result.error) {
        toast.error(result.error);
      } else {
        toast.success("Geofence created successfully");
        addGeofence(result);
        setFormData({
          name: "",
          description: "",
          category: "delivery_zone",
          coordinates: [],
        });
        if (onSuccess) onSuccess();
      }
    } catch (error) {
      toast.error("Failed to create geofence");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6 mb-6">
      <h2 className="text-2xl font-bold mb-4">Create Geofence</h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-1">Name *</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) =>
                setFormData({ ...formData, name: e.target.value })
              }
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="Geofence name"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Category *</label>
            <select
              value={formData.category}
              onChange={(e) =>
                setFormData({ ...formData, category: e.target.value })
              }
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
            >
              <option value="delivery_zone">Delivery Zone</option>
              <option value="restricted_zone">Restricted Zone</option>
              <option value="toll_zone">Toll Zone</option>
              <option value="customer_area">Customer Area</option>
            </select>
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">Description</label>
          <textarea
            value={formData.description}
            onChange={(e) =>
              setFormData({ ...formData, description: e.target.value })
            }
            className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
            placeholder="Geofence description"
            rows="3"
          />
        </div>

        <div className="border rounded-lg p-4 bg-gray-50">
          <h3 className="font-semibold mb-3">Coordinates</h3>

          <div className="grid grid-cols-2 gap-2 mb-3">
            <div>
              <label className="block text-sm mb-1">Latitude</label>
              <input
                id="lat"
                type="number"
                step="0.0001"
                min="-90"
                max="90"
                placeholder="-90 to 90"
                className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>

            <div>
              <label className="block text-sm mb-1">Longitude</label>
              <input
                id="lng"
                type="number"
                step="0.0001"
                min="-180"
                max="180"
                placeholder="-180 to 180"
                className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
          </div>

          <button
            type="button"
            onClick={handleAddCoordinate}
            className="w-full bg-secondary text-white py-2 rounded hover:bg-opacity-90 transition mb-3"
          >
            + Add Coordinate
          </button>

          <div className="space-y-2 max-h-48 overflow-y-auto">
            {formData.coordinates.map((coord, index) => (
              <div
                key={index}
                className="flex justify-between items-center bg-white p-2 rounded"
              >
                <span className="text-sm">
                  {index + 1}. [{coord[0].toFixed(4)}, {coord[1].toFixed(4)}]
                </span>
                <button
                  type="button"
                  onClick={() => handleRemoveCoordinate(index)}
                  className="text-danger hover:text-opacity-70"
                >
                  Remove
                </button>
              </div>
            ))}
          </div>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full bg-primary text-white py-2 rounded hover:bg-opacity-90 transition disabled:opacity-50"
        >
          {loading ? "Creating..." : "Create Geofence"}
        </button>
      </form>
    </div>
  );
};

export default GeofenceForm;
