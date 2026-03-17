import React, { useState, useEffect } from "react";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import Navbar from "./components/Navbar";
import DashboardPage from "./pages/Dashboard";
import GeofencesPage from "./pages/Geofences";
import VehiclesPage from "./pages/Vehicles";
import AlertsPage from "./pages/Alerts";
import { useWebSocket } from "./hooks/useWebSocket";
import "./styles/globals.css";

function App() {
  const [activeTab, setActiveTab] = useState("dashboard");

  // Initialize WebSocket
  useWebSocket();

  const renderPage = () => {
    switch (activeTab) {
      case "dashboard":
        return <DashboardPage />;
      case "geofences":
        return <GeofencesPage />;
      case "vehicles":
        return <VehiclesPage />;
      case "alerts":
        return <AlertsPage />;
      default:
        return <DashboardPage />;
    }
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <Navbar activeTab={activeTab} setActiveTab={setActiveTab} />

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {renderPage()}
      </main>

      <ToastContainer
        position="bottom-right"
        autoClose={5000}
        hideProgressBar={false}
        newestOnTop={false}
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
      />
    </div>
  );
}

export default App;
