import React, { useEffect, useState } from "react";
import {
  MapContainer,
  TileLayer,
  Polygon,
  CircleMarker,
  Popup,
} from "react-leaflet";
import "leaflet/dist/leaflet.css";

const MapComponent = ({
  geofences = [],
  vehicles = {},
  onMapClick = null,
  selectedLocation = null,
}) => {
  const [map, setMap] = useState(null);

  // Get center of map based on geofences
  const getMapCenter = () => {
    if (selectedLocation) {
      return [selectedLocation.lat, selectedLocation.lng];
    }
    if (geofences.length > 0) {
      const coords = geofences[0].coordinates;
      if (coords && coords.length > 0) {
        return [coords[0][0], coords[0][1]];
      }
    }
    return [37.7749, -122.4194]; // Default to San Francisco
  };

  useEffect(() => {
    if (map) {
      map.on("click", (event) => {
        if (onMapClick) {
          onMapClick({
            lat: event.latlng.lat,
            lng: event.latlng.lng,
          });
        }
      });
    }
  }, [map, onMapClick]);

  return (
    <MapContainer
      center={getMapCenter()}
      zoom={13}
      scrollWheelZoom={false}
      style={{ height: "100%", width: "100%" }}
      whenCreated={setMap}
    >
      <TileLayer
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />

      {/* Render geofences */}
      {geofences.map((geofence) => (
        <Polygon
          key={geofence.id}
          positions={geofence.coordinates.map((coord) => [coord[0], coord[1]])}
          color={geofence.category === "restricted_zone" ? "red" : "blue"}
          fillOpacity={0.2}
        >
          <Popup>
            <strong>{geofence.name}</strong>
            <br />
            {geofence.description}
            <br />
            <small>{geofence.category}</small>
          </Popup>
        </Polygon>
      ))}

      {/* Render vehicle markers */}
      {Object.entries(vehicles).map(([vehicleId, location]) => (
        <CircleMarker
          key={vehicleId}
          center={[location.latitude, location.longitude]}
          radius={8}
          fillColor="green"
          color="darkgreen"
          weight={2}
          opacity={1}
          fillOpacity={0.8}
        >
          <Popup>
            <strong>{vehicleId}</strong>
            <br />
            Lat: {location.latitude.toFixed(4)}
            <br />
            Lng: {location.longitude.toFixed(4)}
          </Popup>
        </CircleMarker>
      ))}

      {/* Render selected location */}
      {selectedLocation && (
        <CircleMarker
          center={[selectedLocation.lat, selectedLocation.lng]}
          radius={10}
          fillColor="orange"
          color="darkorange"
          weight={2}
          opacity={1}
          fillOpacity={0.8}
        >
          <Popup>Selected Location</Popup>
        </CircleMarker>
      )}
    </MapContainer>
  );
};

export default MapComponent;
