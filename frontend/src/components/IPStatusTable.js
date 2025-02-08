import React, { useEffect, useState } from 'react';

const IPStatusTable = () => {
  const [ipStatusData, setIpStatusData] = useState([]);

  useEffect(() => {
    // Fetch IP status data from the backend
    fetch("/api/ip")
      .then(response => response.json())
      .then(data => setIpStatusData(data))
      .catch(error => console.error('Error fetching IP status data:', error));
  }, []);

  return (
    <div>
      <h2>IP Status Table</h2>
      <table>
        <thead>
        <tr>
          <th>IP Address</th>
          <th>Status</th>
        </tr>
        </thead>
        <tbody>
        {ipStatusData.map((item, index) => (
          <tr key={index}>
            <td>{item.ip}</td>
            <td>{item.ping_time}</td>
          </tr>
        ))}
      </tbody>
    </table>
</div>
)
  ;
}

export default IPStatusTable;
