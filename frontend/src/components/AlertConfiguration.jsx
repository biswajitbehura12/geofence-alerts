import { useEffect, useState } from "react";
import {
  useVehicleStore,
  useGeofenceStore,
  useAlertStore,
} from "../services/store";
import { useAlertApi } from "../hooks/useApi";
import { toast } from "react-toastify";

const AlertConfiguration = () => {
  const { vehicles = [] } = useVehicleStore();
  const { geofences = [] } = useGeofenceStore();
  const { alerts = [], setAlerts } = useAlertStore();
  const { configureAlert, getAlerts } = useAlertApi();
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    geofence_id: "",
    vehicle_id: "",
    event_type: "entry",
  });

  useEffect(() => {
    loadAlerts();
  }, []);

  const loadAlerts = async () => {
    try {
      const result = await getAlerts();
      if (result.alerts) {
        setAlerts(result.alerts);
      }
    } catch (error) {
      console.error("Failed to load alerts", error);
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!formData.geofence_id) {
      toast.error("Geofence is required");
      return;
    }

    setLoading(true);
    try {
      const result = await configureAlert(formData);

      if (result.error) {
        toast.error(result.error);
      } else {
        toast.success("Alert configured successfully");
        await loadAlerts();
        setFormData({
          geofence_id: "",
          vehicle_id: "",
          event_type: "entry",
        });
      }
    } catch (error) {
      toast.error("Failed to configure alert");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-2xl font-bold mb-4">Configure Alerts</h2>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">Geofence *</label>
            <select
              name="geofence_id"
              value={formData.geofence_id}
              onChange={handleChange}
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
            >
              <option value="">Select a geofence</option>
              {geofences.map((g) => (
                <option key={g.id} value={g.id}>
                  {g.name} ({g.category})
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">
              Vehicle (optional - leave empty for all vehicles)
            </label>
            <select
              name="vehicle_id"
              value={formData.vehicle_id}
              onChange={handleChange}
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
            >
              <option value="">All vehicles</option>
              {vehicles.map((v) => (
                <option key={v.id} value={v.id}>
                  {v.vehicle_number} - {v.driver_name}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">
              Event Type *
            </label>
            <select
              name="event_type"
              value={formData.event_type}
              onChange={handleChange}
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
            >
              <option value="entry">Entry</option>
              <option value="exit">Exit</option>
              <option value="both">Both</option>
            </select>
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-primary text-white py-2 rounded hover:bg-opacity-90 transition disabled:opacity-50"
          >
            {loading ? "Configuring..." : "Configure Alert"}
          </button>
        </form>
      </div>

      {/* Configured Alerts */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h3 className="text-xl font-bold mb-4">Configured Alerts</h3>

        {alerts.length === 0 ? (
          <p className="text-gray-500 text-center py-8">
            No alerts configured yet
          </p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b">
                  <th className="text-left py-2 px-2">Geofence</th>
                  <th className="text-left py-2 px-2">Vehicle</th>
                  <th className="text-left py-2 px-2">Event</th>
                  <th className="text-left py-2 px-2">Status</th>
                  <th className="text-left py-2 px-2">Created</th>
                </tr>
              </thead>
              <tbody>
                {alerts.map((alert) => (
                  <tr
                    key={alert.alert_id}
                    className="border-b hover:bg-gray-50"
                  >
                    <td className="py-2 px-2">{alert.geofence_name}</td>
                    <td className="py-2 px-2">
                      {alert.vehicle_number || "All"}
                    </td>
                    <td className="py-2 px-2">{alert.event_type}</td>
                    <td className="py-2 px-2">
                      <span className="bg-success text-white px-2 py-1 rounded text-xs">
                        {alert.status}
                      </span>
                    </td>
                    <td className="py-2 px-2">
                      {new Date(alert.created_at).toLocaleDateString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
};

export default AlertConfiguration;
