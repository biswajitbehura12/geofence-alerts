import { useCallback, useEffect, useState } from "react";
import {
  useAlertStore,
  useVehicleStore,
  useGeofenceStore,
} from "../services/store";
import { useAlertApi } from "../hooks/useApi";

const ViolationHistory = () => {
  const { violations = [], setViolations } = useAlertStore();
  const { vehicles = [] } = useVehicleStore();
  const { geofences = [] } = useGeofenceStore();
  const { getViolationHistory } = useAlertApi();
  const [loading, setLoading] = useState(false);
  const [filters, setFilters] = useState({
    vehicle_id: "",
    geofence_id: "",
    start_date: "",
    end_date: "",
    limit: 50,
    offset: 0,
  });

  const loadViolationHistory = useCallback(async () => {
    setLoading(true);
    try {
      const filterObj = {};
      if (filters.vehicle_id) filterObj.vehicleId = filters.vehicle_id;
      if (filters.geofence_id) filterObj.geofenceId = filters.geofence_id;
      if (filters.start_date)
        filterObj.startDate = new Date(filters.start_date).toISOString();
      if (filters.end_date)
        filterObj.endDate = new Date(filters.end_date).toISOString();
      filterObj.limit = filters.limit;
      filterObj.offset = filters.offset;

      const result = await getViolationHistory(filterObj);
      if (result.violations) {
        setViolations(result.violations);
      }
    } catch (error) {
      console.error("Failed to load violations", error);
    } finally {
      setLoading(false);
    }
  }, [filters, getViolationHistory, setViolations]);

  useEffect(() => {
    loadViolationHistory();
  }, [loadViolationHistory]);

  const handleFilterChange = (e) => {
    const { name, value } = e.target;
    setFilters({ ...filters, [name]: value, offset: 0 });
  };

  const handleApplyFilters = () => {
    loadViolationHistory();
  };

  const getEventColor = (eventType) => {
    return eventType === "entry" ? "bg-warning" : "bg-danger";
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-2xl font-bold mb-4">Violation History</h2>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-3 mb-6 p-4 bg-gray-50 rounded-lg">
        <div>
          <label className="block text-sm font-medium mb-1">Vehicle</label>
          <select
            name="vehicle_id"
            value={filters.vehicle_id}
            onChange={handleFilterChange}
            className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-primary"
          >
            <option value="">All</option>
            {vehicles.map((v) => (
              <option key={v.id} value={v.id}>
                {v.vehicle_number}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">Geofence</label>
          <select
            name="geofence_id"
            value={filters.geofence_id}
            onChange={handleFilterChange}
            className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-primary"
          >
            <option value="">All</option>
            {geofences.map((g) => (
              <option key={g.id} value={g.id}>
                {g.name}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">From</label>
          <input
            type="date"
            name="start_date"
            value={filters.start_date}
            onChange={handleFilterChange}
            className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-primary"
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">To</label>
          <input
            type="date"
            name="end_date"
            value={filters.end_date}
            onChange={handleFilterChange}
            className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-primary"
          />
        </div>
      </div>

      <button
        onClick={handleApplyFilters}
        disabled={loading}
        className="mb-4 bg-primary text-white px-6 py-2 rounded hover:bg-opacity-90 transition disabled:opacity-50"
      >
        {loading ? "Loading..." : "Apply Filters"}
      </button>

      {violations.length === 0 ? (
        <p className="text-gray-500 text-center py-8">No violations found</p>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b">
                <th className="text-left py-3 px-3">Vehicle</th>
                <th className="text-left py-3 px-3">Geofence</th>
                <th className="text-left py-3 px-3">Event</th>
                <th className="text-left py-3 px-3">Location</th>
                <th className="text-left py-3 px-3">Timestamp</th>
              </tr>
            </thead>
            <tbody>
              {violations.map((violation) => (
                <tr key={violation.id} className="border-b hover:bg-gray-50">
                  <td className="py-3 px-3 font-medium">
                    {violation.vehicle_number}
                  </td>
                  <td className="py-3 px-3">{violation.geofence_name}</td>
                  <td className="py-3 px-3">
                    <span
                      className={`${getEventColor(violation.event_type)} text-white px-2 py-1 rounded text-xs`}
                    >
                      {violation.event_type.toUpperCase()}
                    </span>
                  </td>
                  <td className="py-3 px-3 text-xs">
                    {violation.latitude.toFixed(4)},{" "}
                    {violation.longitude.toFixed(4)}
                  </td>
                  <td className="py-3 px-3 text-xs">
                    {new Date(violation.timestamp).toLocaleString()}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default ViolationHistory;
