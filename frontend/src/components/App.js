import React, { useEffect, useState } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';
import IPStatusTable from './IPStatusTable';

const App = () => {
  const [ipStatuses, setIpStatuses] = useState([]);

  // Function to fetch IP statuses from the backend
  const fetchIPStatuses = async () => {
    try {
      const response = await fetch('/api/ip');
      if (response.ok) {
        const data = await response.json();
        setIpStatuses(data);
      } else {
        console.error('Failed to fetch IP statuses');
      }
    } catch (error) {
      console.error('Error fetching IP statuses:', error);
    }
  };

  useEffect(() => {
    fetchIPStatuses();  // Initial fetch
    const interval = setInterval(fetchIPStatuses, 5000); // Re-fetch every 5 seconds

    return () => clearInterval(interval); // Clear interval on cleanup
  }, []);

  return (
    <div className="container mt-4">
      <h1 className="text-center text-primary mb-4">IP Statuses</h1>
      <IPStatusTable ipStatuses={ipStatuses} />
    </div>
  );
};

export default App;
