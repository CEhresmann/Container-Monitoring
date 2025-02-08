import React from 'react';
import { Table, Dropdown } from 'react-bootstrap';

const IPStatusTable = ({ ipStatuses }) => {
  const groupedStatuses = ipStatuses.reduce((acc, status) => {
    if (!acc[status.ip]) acc[status.ip] = [];
    acc[status.ip].push(status);
    return acc;
  }, {});

  return (
    <div className="container mt-4">
      <Table striped bordered hover responsive variant="light" className="shadow-sm">
        <thead className="bg-primary text-white">
        <tr>
          <th>IP адрес</th>
          <th>Время пинга (в ms)</th>
          <th>Дата последней успешной попытки</th>
        </tr>
        </thead>
        <tbody>
        {Object.keys(groupedStatuses).map((ip, index) => (
          <tr key={index} className={index % 2 === 0 ? 'table-secondary' : ''}>
            <td>
              {/* Выпадающий список с историей пинга для каждого IP */}
              <Dropdown>
                <Dropdown.Toggle variant="success" id={`dropdown-${ip}`}>
                  {ip}
                </Dropdown.Toggle>
                <Dropdown.Menu>
                  {groupedStatuses[ip].length > 0 ? (
                    groupedStatuses[ip].map((status, idx) => (
                      <Dropdown.Item key={idx}>
                        Ping: {status.ping_time} ms, {new Date(status.last_ok).toLocaleString()}
                      </Dropdown.Item>
                    ))
                  ) : (
                    <Dropdown.Item>Нет данных</Dropdown.Item>
                  )}
                </Dropdown.Menu>
              </Dropdown>
            </td>
            <td>{groupedStatuses[ip][groupedStatuses[ip].length - 1].ping_time}</td>
            <td>{new Date(groupedStatuses[ip][groupedStatuses[ip].length - 1].last_ok).toLocaleString()}</td>
          </tr>
        ))}
        </tbody>
      </Table>
    </div>
  );
};

export default IPStatusTable;
