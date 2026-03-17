import React, { useEffect, useState } from "react";
import { useAlertStore } from "../services/store";

const AlertsFeed = () => {
  const { recentAlerts = [] } = useAlertStore();
  const [expandedId, setExpandedId] = useState(null);

  const getCategoryColor = (category) => {
    switch (category) {
      case "restricted_zone":
        return "bg-danger text-white";
      case "delivery_zone":
        return "bg-success text-white";
      case "toll_zone":
        return "bg-warning text-black";
      default:
        return "bg-secondary text-white";
    }
  };

  const getEventIcon = (eventType) => {
    return eventType === "entry" ? "📍" : "🚫";
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-2xl font-bold mb-4">Real-time Alerts</h2>

      {recentAlerts.length === 0 ? (
        <p className="text-gray-500 text-center py-8">No alerts yet</p>
      ) : (
        <div className="space-y-3 max-h-96 overflow-y-auto">
          {recentAlerts.map((alert) => (
            <div
              key={alert.event_id}
              className={`p-4 rounded-lg border-l-4 ${getCategoryColor(alert.geofence.category)}`}
            >
              <div className="flex justify-between items-start">
                <div>
                  <div className="flex items-center gap-2 mb-2">
                    <span className="text-2xl">
                      {getEventIcon(alert.event_type)}
                    </span>
                    <strong>{alert.vehicle.vehicle_number}</strong>
                    <span className="uppercase text-sm font-semibold">
                      {alert.event_type}
                    </span>
                  </div>
                  <p className="text-sm">
                    <strong>{alert.geofence.geofence_name}</strong>
                  </p>
                  <p className="text-xs opacity-75">
                    Driver: {alert.vehicle.driver_name}
                  </p>
                  <p className="text-xs opacity-60">
                    {new Date(alert.timestamp).toLocaleString()}
                  </p>
                </div>

                <button
                  onClick={() =>
                    setExpandedId(
                      expandedId === alert.event_id ? null : alert.event_id,
                    )
                  }
                  className="text-xl"
                >
                  {expandedId === alert.event_id ? "▼" : "▶"}
                </button>
              </div>

              {expandedId === alert.event_id && (
                <div className="mt-3 pt-3 border-t opacity-75 text-sm">
                  <p>
                    Location: {alert.location.latitude.toFixed(4)},{" "}
                    {alert.location.longitude.toFixed(4)}
                  </p>
                  <p>Category: {alert.geofence.category}</p>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default AlertsFeed;
