import React, { useEffect } from "react";
import AlertConfiguration from "../components/AlertConfiguration";
import ViolationHistory from "../components/ViolationHistory";
import AlertsFeed from "../components/AlertsFeed";

const AlertsPage = () => {
  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold">Alert Management</h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Alerts Feed */}
        <div className="lg:col-span-1">
          <AlertsFeed />
        </div>

        {/* Configuration */}
        <div className="lg:col-span-2">
          <AlertConfiguration />
        </div>
      </div>

      {/* Violation History */}
      <ViolationHistory />
    </div>
  );
};

export default AlertsPage;
