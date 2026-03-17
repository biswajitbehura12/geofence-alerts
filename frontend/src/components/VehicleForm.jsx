import { useState } from "react";
import { useVehicleStore } from "../services/store";
import { useVehicleApi } from "../hooks/useApi";
import { toast } from "react-toastify";

const VehicleForm = ({ onSuccess }) => {
  const { addVehicle } = useVehicleStore();
  const { registerVehicle } = useVehicleApi();
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    vehicle_number: "",
    driver_name: "",
    vehicle_type: "truck",
    phone: "",
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!formData.vehicle_number || !formData.driver_name || !formData.phone) {
      toast.error("All fields are required");
      return;
    }

    setLoading(true);
    try {
      const result = await registerVehicle(formData);

      if (result.error) {
        toast.error(result.error);
      } else {
        toast.success("Vehicle registered successfully");
        addVehicle(result);
        setFormData({
          vehicle_number: "",
          driver_name: "",
          vehicle_type: "truck",
          phone: "",
        });
        if (onSuccess) onSuccess();
      }
    } catch (error) {
      toast.error("Failed to register vehicle");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6 mb-6">
      <h2 className="text-2xl font-bold mb-4">Register Vehicle</h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-1">
              Vehicle Number *
            </label>
            <input
              type="text"
              name="vehicle_number"
              value={formData.vehicle_number}
              onChange={handleChange}
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="KA-01-AB-1234"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">
              Driver Name *
            </label>
            <input
              type="text"
              name="driver_name"
              value={formData.driver_name}
              onChange={handleChange}
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="John Doe"
            />
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-1">
              Vehicle Type *
            </label>
            <select
              name="vehicle_type"
              value={formData.vehicle_type}
              onChange={handleChange}
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
            >
              <option value="truck">Truck</option>
              <option value="car">Car</option>
              <option value="van">Van</option>
              <option value="motorcycle">Motorcycle</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Phone *</label>
            <input
              type="tel"
              name="phone"
              value={formData.phone}
              onChange={handleChange}
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="+1234567890"
            />
          </div>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full bg-primary text-white py-2 rounded hover:bg-opacity-90 transition disabled:opacity-50"
        >
          {loading ? "Registering..." : "Register Vehicle"}
        </button>
      </form>
    </div>
  );
};

export default VehicleForm;
