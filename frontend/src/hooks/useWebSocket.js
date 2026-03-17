import { useEffect, useState } from "react";
import { useAlertStore, useWebSocketStore } from "../services/store";
import { toast } from "react-toastify";
import { createWebSocketConnection } from "../services/api";

export const useWebSocket = () => {
  const { setConnected } = useWebSocketStore();
  const { addAlert } = useAlertStore();
  const [ws, setWs] = useState(null);

  useEffect(() => {
    try {
      const webSocket = createWebSocketConnection();

      webSocket.onopen = () => {
        console.log("WebSocket connected");
        setConnected(true);
        toast.success("Connected to real-time alerts");
      };

      webSocket.onmessage = (event) => {
        const alert = JSON.parse(event.data);
        addAlert(alert);

        // Show toast notification
        const message = `Vehicle ${alert.vehicle.vehicle_number} ${alert.event_type} ${alert.geofence.geofence_name}`;
        if (alert.geofence.category === "restricted_zone") {
          toast.error(message, { autoClose: false });
        } else {
          toast.info(message);
        }
      };

      webSocket.onerror = (error) => {
        console.error("WebSocket error:", error);
        toast.error("WebSocket connection error");
      };

      webSocket.onclose = () => {
        console.log("WebSocket disconnected");
        setConnected(false);
      };

      setWs(webSocket);
    } catch (error) {
      console.error("Failed to connect WebSocket:", error);
      toast.error("Failed to connect to real-time alerts");
    }

    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, [setConnected, addAlert]);

  return ws;
};
