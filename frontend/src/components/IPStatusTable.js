import React, { useEffect, useState } from 'react';
import { Table } from 'react-bootstrap';

function IPStatusTable() {
    const [ipStatuses, setIpStatuses] = useState([]);

    useEffect(() => {
        fetch('/api/ips')
            .then((response) => response.json())
            .then((data) => setIpStatuses(data))
            .catch((error) => console.error('Error fetching IP statuses:', error));
    }, []);

    return (
        <div>
            <h1>IP Statuses</h1>
            <Table striped bordered hover>
                <thead>
                <tr>
                    <th>IP Address</th>
                    <th>Ping Time</th>
                    <th>Last OK</th>
                </tr>
                </thead>
                <tbody>
                {ipStatuses.map((status, index) => (
                    <tr key={index}>
                        <td>{status.ip}</td>
                        <td>{status.ping_time}</td>
                        <td>{status.last_ok}</td>
                    </tr>
                ))}
                </tbody>
            </Table>
        </div>
    );
}

export default IPStatusTable;
