import logo from '../logo.svg';
import '../App.css';
import React, { useEffect } from 'react';
import IPStatusTable from './IPStatusTable';

function App() {
    useEffect(() => {
        fetch('/api/ip')
            .then((response) => response.json())
            .then((data) => {
                console.log('IP статусы:', data);
            })
            .catch((error) => console.error('Ошибка при получении IP статусов:', error));
    }, []);

    const addIPStatus = (ip, pingTime, lastOK) => {
        fetch('/api/ip', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ ip, pingTime, lastOK }),
        })
            .then((response) => {
                if (response.ok) {
                    console.log('IP статус добавлен');
                }
            })
            .catch((error) => console.error('Ошибка при добавлении IP статуса:', error));
    };

    return (
        <div className="container">
            <IPStatusTable />
            {/* Здесь можно добавить элементы UI, связанные с addIPStatus */}
        </div>
    );
}

export default App;
