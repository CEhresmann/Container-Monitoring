import React from 'react';
import { Table } from 'react-bootstrap';

const IPStatusTable = ({ ipStatuses }) => {
  return (
    <div className="container mt-4">
      <h2 className="text-center text-primary mb-4">IP Status Table</h2>
      <Table striped bordered hover responsive variant="light" className="shadow-sm">
        <thead className="bg-primary text-white">
        <tr>
          <th>IP адрес</th>
          <th>Время пинга (в ms)</th>
          <th>Дата последней успешной попытки</th>
        </tr>
        </thead>
        <tbody>
        {ipStatuses.map((status, index) => (
          <tr key={index} className={index % 2 === 0 ? 'table-secondary' : ''}>
            <td>{status.ip}</td>
            <td>{status.ping_time}</td>
            <td>{new Date(status.last_ok).toLocaleString()}</td>
          </tr>
        ))}
        </tbody>
      </Table>
    </div>
  );
};

export default IPStatusTable;
