import React, { useEffect, useState } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';
import IPStatusTable from './IPStatusTable';

const App = () => {
  const [ipStatuses, setIpStatuses] = useState([]);
  const [loading, setLoading] = useState(true);

  const fetchIPStatuses = async () => {
    try {
      const response = await fetch('/api/ip');
      if (response.ok) {
        const data = await response.json();
        setIpStatuses(data);
        setLoading(false);
      } else {
        console.error('Failed to fetch IP statuses');
        setLoading(false);
      }
    } catch (error) {
      console.error('Error fetching IP statuses:', error);
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchIPStatuses();
    const interval = setInterval(fetchIPStatuses, 5000);

    return () => clearInterval(interval);
  }, []);

  return (
    <div className="container mt-4" style={{ background: 'url("/path-to-your-background-image.jpg")', backgroundSize: 'cover', backgroundPosition: 'center', backdropFilter: 'blur(5px)' }}>
      <h1 className="text-center text-primary mb-4">IP Statuses</h1>
      {loading ? (
        <div className="d-flex justify-content-center align-items-center" style={{ minHeight: '400px' }}>
          <div className="spinner-border text-primary" role="status">
            <span className="visually-hidden">Loading...</span>
          </div>
        </div>
      ) : (
        <IPStatusTable ipStatuses={ipStatuses} />
      )}
    </div>
  );
};

export default App;
